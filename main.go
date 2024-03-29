package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	applicationName   = "prometheus_export"
	letterBytes       = "abcdefghijklmnopqrstuvwxyz" // map for random string
	letterIdxBits     = 6                            // 6 bits to represent a letter index
	letterIdxMask     = 1<<letterIdxBits - 1         // All 1-bits, as many as letterIdxBits
	letterIdxMax      = 63 / letterIdxBits           // # of letter indices fitting in 63 bits
	timeRangeOverSize = 15                           // minutes used for prolong data before and after Str/Stop
)

var (
	logger        log.Logger // logger
	Version       string
	Revision      string
	Branch        string
	BuildUser     string
	BuildDate     string
	src           = rand.NewSource(time.Now().UnixNano()) // randomize base string
	timeStart     time.Time
	maxRandomSize = 10 // required size of random string
	//metricsMeta   *MetricsMetaList
)

func RandomString() string {
	sb := strings.Builder{}
	sb.Grow(maxRandomSize)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := maxRandomSize-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	version.Branch = Branch
	version.Revision = Revision
	version.BuildUser = BuildUser
	version.BuildDate = BuildDate
	version.Version = Version
	kingpin.Version(version.Print(applicationName))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger = promlog.New(promlogConfig)
	_ = level.Info(logger).Log("msg", "Starting prometheus data export ", "version", version.Info())
	timeStart = time.Now().UTC()

	err := config.LoadFile(*configFile)

	if *showConfig {
		_ = level.Info(logger).Log("msg", "show only configuration ane exit")

		if err != nil {
			fmt.Print(config.print())
		} else {
			if config.useRange() {
				fmt.Printf("%s%s", config.print(), printTimeRanges(generateRangeTable(initRangeFromTo(configFrom, configTo))))
			} else {
				fmt.Printf("%s%s", config.print(), printTimeRanges(generateRangeTable(initRange())))
			}
		}
		os.Exit(0)
	}
	if err != nil {
		_ = level.Error(logger).Log("msg", "problem with configuration", "error", err)
		fmt.Printf("Program did not start due to configuration error! \r\n\tError: %s", err)
		os.Exit(1)
	}

	// start process data
	_ = level.Info(logger).Log("msg", "start collect necessary data")
	v1Api, err := prepareApi(config)
	if err != nil {
		_ = level.Info(logger).Log("msg", "program exit, because there is problem with connect Prometheus server")
		os.Exit(2)
	}

	start := time.Now()
	var dateRange v1.Range
	if config.useRange() {
		dateRange = initRangeFromTo(configFrom, configTo)
	} else {
		dateRange = initRange()
	}
	timeRangeSplit := generateRangeTable(dateRange)
	_ = level.Info(logger).Log("msg", fmt.Sprintf("requried time split to %d time ranges", len(timeRangeSplit)))

	// process for all data
	services := OrganizedServices{}
	for _, source := range config.Sources {
		services = processOneSources(v1Api, source, services, dateRange)
	}
	// process labels
	if len(config.Labels) > 0 {
		services = processLabels(v1Api, config.Labels, services, dateRange)
	}
	startReadApi(v1Api, timeRangeSplit, services)
	_ = level.Info(logger).Log("msg", fmt.Sprintf("program successful ends after %s", time.Since(start).String()))
}

func processOneSources(api v1.API, sources Sources, services OrganizedServices, timeRange v1.Range) OrganizedServices {
	labels, err := collectSeriesList(api, sources, timeRange)
	if err != nil {
		return services
	}

	storage := NewStorage(config.Path, sources)
	_ = storage.prepareDirectory()
	organized := splitToSeriesNameAndInstance(labels, storage)
	storage.saveOrganized(organized, timeRange)

	return append(services, organized...)
}

func processLabels(api v1.API, lbl []Labels, services OrganizedServices, timeRange v1.Range) OrganizedServices {
	labels, err := collectLabelsSeriesList(api, lbl, timeRange)
	if err != nil {
		return services
	}

	storage := NewNameStorage(config.Path, "labels")
	_ = storage.prepareDirectory()
	organized := splitToSeriesNameAndInstance(labels, storage)
	storage.saveOrganized(organized, timeRange)

	return append(services, organized...)
}

func startReadApi(api v1.API, ranges []v1.Range, organized OrganizedServices) {
	channel := make(chan OrganizedSeries, len(organized))

	// fill data
	for _, series := range organized {
		channel <- series
	}
	_ = level.Info(logger).Log("msg", "prepared GO coroutines channel data")

	var wg sync.WaitGroup
	cpu := runtime.NumCPU()
	if cpu > 1 {
		cpu--
	}
	start := time.Now()
	for i := 0; i < cpu; i++ {
		wg.Add(1)
		go processOneInstance(&wg, channel, api, ranges, i)
	}
	_ = level.Info(logger).Log("msg", fmt.Sprintf("wait to finish all %d routines", cpu))
	wg.Wait()
	_ = level.Info(logger).Log("msg", fmt.Sprintf("all coroutines finish %s", time.Since(start).String()))
}

func processOneInstance(wg *sync.WaitGroup, channel chan OrganizedSeries, api v1.API, ranges []v1.Range, cpu int) {
	defer wg.Done()
	for {
		select {
		case series := <-channel:
			_ = level.Debug(logger).Log("msg", fmt.Sprintf("process data for coroutine %d", cpu), "series", series.Name)
			var saveAll []SaveAllData
			for _, label := range series.Set {
				var saveData []SaveData
				loops := 0
				loopsErr := 0
				loopsEmpty := 0
				for _, timeRange := range ranges {
					loops++
					data, err := readQueryRange(api, label, timeRange)
					if err != nil {
						loopsErr++
						continue
					}
					d := convertApiData(data)
					if len(d) > 0 {
						saveData = append(saveData, d...)
					} else {
						loopsEmpty++
					}
				}
				saveAll = append(saveAll, SaveAllData{Collection: label, Data: saveData})
			}
			series.Storage.saveAllData(saveAll, cleanFilePathName(series.Name)+".json")
			break
		case <-time.Tick(1 * time.Second):
			_ = level.Debug(logger).Log("msg", fmt.Sprintf("funish coroutine %d", cpu))
			return
		}
	}
}
