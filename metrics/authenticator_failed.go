package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"regexp"
	"strings"
)

type authenticatorFailedMetric struct {
	counter *prometheus.CounterVec
}

func NewAuthenticatorFailedMetric() *authenticatorFailedMetric {
	return &authenticatorFailedMetric{
		counter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "exim_authenticator_failed_total",
			Help: "Total number of rejected authentication attempts.",
		}, []string{
			"authenticator", "error_message",
		}),
	}
}

func (m *authenticatorFailedMetric) Collector() prometheus.Collector {
	return m.counter
}

func (m *authenticatorFailedMetric) Observe(line string) {
	if strings.Contains(line, "authenticator failed") {
		authenticator := m.parseAuthenticator(line)
		message := m.parseErrorMessage(line)
		m.counter.WithLabelValues(authenticator, message).Inc()
	}
}

func (m *authenticatorFailedMetric) parseAuthenticator(line string) string {
	r := regexp.MustCompile("([a-zA-Z][^\\s]+) authenticator failed")
	matches := r.FindStringSubmatch(line)
	if len(matches) > 1 && len(matches[1]) > 0 {
		return matches[1]
	}
	return "unknown authenticator"
}

func (m *authenticatorFailedMetric) parseErrorMessage(line string) string {
	r := regexp.MustCompile("authenticator failed.*:\\s*([^\\(]+)")
	matches := r.FindStringSubmatch(line)
	if len(matches) > 1 && len(matches[1]) > 0 {
		return strings.TrimSpace(matches[1])
	}
	return "unknown error message"
}
