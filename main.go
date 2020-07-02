package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	applicationName = "prometheus_export"
	letterBytes     = "abcdefghijklmnopqrstuvwxyz" // map for random string
	letterIdxBits   = 6                            // 6 bits to represent a letter index
	letterIdxMask   = 1<<letterIdxBits - 1         // All 1-bits, as many as letterIdxBits
	letterIdxMax    = 63 / letterIdxBits           // # of letter indices fitting in 63 bits
)

var (
	logger        log.Logger // logger
	Version       string
	Revision      string
	Branch        string
	BuildUser     string
	BuildDate     string
	src           = rand.NewSource(time.Now().UnixNano()) // randomize base string
	maxRandomSize = 10                                    // required size of random string
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

func containsString(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func saveMetaData() *MetricsMetaList {
	l, _ := readTargetsList()
	metricsMeta := NewMetricsMetaList(*l)
	metricsMeta.onlyForJobs(config.Jobs)
	metricsMeta.saveList()

	_ = level.Info(logger).Log("msg", fmt.Sprintf("metrics in meta data %d", len(*metricsMeta)))
	return metricsMeta
}

func saveData(metricsMeta *MetricsMetaList) {
	var wg sync.WaitGroup
	_ = level.Debug(logger).Log("msg", fmt.Sprintf("start %d goroutines for collect metrics for %d days ", len(*metricsMeta), config.Days))
	for _, metrics := range *metricsMeta {
		wg.Add(1)
		go collectOneMetrics(&wg, metrics.Metric)
	}
	_ = level.Debug(logger).Log("msg", "all goroutines started")
	wg.Wait()
	_ = level.Debug(logger).Log("msg", "data store finish")
}

func collectOneMetrics(wg *sync.WaitGroup, metricName string) {
	defer wg.Done()
	m := Matrix{}
	for i := 0; i < config.Days; i++ {
		c, err := getRangeDay(metricName+"{}", i+1)
		if err != nil {
			continue
		}
		l := len(m)
		for _, series := range *c {
			m.appendSeries(series)
		}
		_ = level.Debug(logger).Log("msg", fmt.Sprintf("in metrics %s add new %d series", metricName, len(m)-l))
	}
	if len(m) == 0 {
		_ = level.Error(logger).Log("msg", "for metrics "+metricName+" not any data")
	} else {
		m.save(metricName)
		_ = level.Debug(logger).Log("msg", "save data for metrics "+metricName)
	}
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

	err := config.LoadFile(*configFile)
	if *showConfig {
		_ = level.Info(logger).Log("msg", "show only configuration ane exit")
		fmt.Print(config.print())
		os.Exit(0)
	}
	if err != nil {
		_ = level.Error(logger).Log("msg", "problem with configuration", "error", err)
		fmt.Printf("Program did not start due to configuration error! \r\n\tError: %s", err)
		os.Exit(1)
	}
	m := saveMetaData()
	saveData(m)
}
