package main

import (
	"flag"
	"fmt"
	"github.com/ActiveState/tail"
	"github.com/fstab/exim_prometheus_exporter/metrics"
	"github.com/fstab/exim_prometheus_exporter/server"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	mainlog = flag.String("mainlog", "/var/log/exim4/mainlog", "Path to the exim4 mainlog file.")
	readall = flag.Bool("readall", false, "Read mainlog from beginning? Default (false) is to start at the end of the file and read only new lines.")
	port    = flag.Int("port", 8443, "The port to listen on for HTTPS requests.")
	cert    = flag.String("cert", "", "Path to the SSL certificate file for the HTTPS server. (optional)")
	key     = flag.String("key", "", "Path to the SSL private key file for the HTTPS server. (optional)")
)

var (
	// total keeps track of how many log lines are processed by the Metrics.
	total = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "exim_loglines_total",
		Help: "Total number of lines in Exim's mainlog. The 'processed' flag tells if exim_prometheus_exporter has processed this line (if false the line was ignored).",
	}, []string{
		"processed", "metric",
	})
)

func main() {
	parseCommandline()
	initPrometheus()
	serverErrorChannel := startServer("/metrics", prometheus.Handler())
	fmt.Printf("Starting server on https://localhost:%v/metrics\n", *port)
	err := processLogLines(serverErrorChannel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(-1)
	}
}

func parseCommandline() {
	flag.Parse()
	if len(*cert) == 0 && len(*key) > 0 || len(*cert) > 0 && len(*key) == 0 {
		fmt.Fprintln(os.Stderr, "Syntax error: -cert and -key cannot be used without each other.")
		os.Exit(-1)
	}
}

func initPrometheus() {
	for _, m := range metrics.Metrics {
		prometheus.MustRegister(m.Collector())
	}
	prometheus.MustRegister(total)
}

func startServer(path string, handler http.Handler) chan error {
	result := make(chan error)
	go func() {
		if len(*cert) > 0 && len(*key) > 0 {
			result <- server.Run(*port, *cert, *key, path, handler)
		} else {
			result <- server.RunWithDefaultKeys(*port, path, handler)
		}
	}()
	return result
}

func processLogLines(serverErrorChannel chan error) error {
	tailFile, err := tailFile(*mainlog, *readall)
	if err != nil {
		return fmt.Errorf("Failed to read %v: %v", *mainlog, err.Error())
	}
	for {
		select {
		case err := <-serverErrorChannel:
			tailFile.Stop()
			return fmt.Errorf("Server error: %v", err.Error())
		case line := <-tailFile.Lines:
			if line.Err == nil {
				process(line.Text)
			}
		}
	}
}

func tailFile(logfile string, readall bool) (*tail.Tail, error) {
	whence := os.SEEK_END
	if readall {
		whence = os.SEEK_SET
	}
	return tail.TailFile(logfile, tail.Config{
		MustExist: true,                      // Fail early if the file does not exist
		ReOpen:    true,                      // Reopen recreated files (tail -F)
		Follow:    true,                      // Continue looking for new lines (tail -f)
		Logger:    tail.DiscardingLogger,     // Disable logging
		Location:  &tail.SeekInfo{0, whence}, // Start at the beginning or end of the file?
	})
}

func process(line string) {
	processed := false
	metricNames := make([]string, 0)
	for _, metric := range metrics.Metrics {
		if metric.Matches(line) {
			metric.Process(line)
			processed = true
			metricNames = append(metricNames, metric.Name())
		}
	}
	total.WithLabelValues(strconv.FormatBool(processed), strings.Join(metricNames, ", ")).Inc()
}
