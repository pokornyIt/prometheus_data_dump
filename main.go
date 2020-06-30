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

	l, _ := readTargetsList()
	fmt.Printf("Targets: %d", len(*l))
}
