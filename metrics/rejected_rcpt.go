package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"strings"
)

const exim_rejected_rcpt_total = "exim_rejected_rcpt_total"

type rejectedRcptMetric struct {
	counter *prometheus.CounterVec
}

func NewRejectedRcptMetric() *rejectedRcptMetric {
	return &rejectedRcptMetric{
		counter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: exim_rejected_rcpt_total,
			Help: "Total number of rejected recipients, partitioned by error message.",
		}, []string{
			"error_message",
		}),
	}
}

func (m *rejectedRcptMetric) Name() string {
	return exim_rejected_rcpt_total
}

func (m *rejectedRcptMetric) Collector() prometheus.Collector {
	return m.counter
}

func (m *rejectedRcptMetric) Matches(line string) bool {
	return strings.Contains(line, "rejected RCPT") && !strings.Contains(line, "temporarily rejected RCPT")
}

func (m *rejectedRcptMetric) Process(line string) {
	msg := m.parseRejectedRcptErrorMessage(line)
	m.counter.WithLabelValues(msg).Inc()
}

func (m *rejectedRcptMetric) parseRejectedRcptErrorMessage(line string) string {
	var result string
	parts := strings.Split(line, ":") // after ':' should be the error message
	if len(parts) >= 2 && strings.Contains(parts[len(parts)-2], "rejected RCPT") {
		result = strings.TrimSpace(parts[len(parts)-1])
	}
	if len(result) == 0 {
		fmt.Fprintf(os.Stderr, "Found 'rejected RCPT' log line with unknown error message: %v", line)
		return "unknown error"
	}
	return result
}
