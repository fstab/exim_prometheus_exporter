package metrics

import "testing"

func TestNoDuplicateNames(t *testing.T) {
	names := make(map[string]bool)
	for _, metric := range Metrics {
		_, found := names[metric.Name()]
		if found {
			t.Fail() // duplicate name
		}
		names[metric.Name()] = true
	}
}
