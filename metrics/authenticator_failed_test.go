package metrics

import (
	"testing"
)

func TestParseAuthenticatorOk(t *testing.T) {
	line := "2016-01-21 16:16:38 cram_md5_server authenticator failed for ([192.168.0.3]) [75.148.200.89]: 535 Incorrect authentication data (set_id=charles)"
	metric := NewAuthenticatorFailedMetric()
	if "cram_md5_server" != metric.parseAuthenticator(line) {
		t.Fail()
	}
}

func TestParseAuthenticatorFailed(t *testing.T) {
	line := "2016-01-21 16:16:38 authenticator failed for ([192.168.0.3]) [75.148.200.89]: 535 Incorrect authentication data (set_id=charles)"
	metric := NewAuthenticatorFailedMetric()
	if "unknown authenticator" != metric.parseAuthenticator(line) {
		t.Fail()
	}
}

func TestParseErrorMessageOk(t *testing.T) {
	line := "2016-01-21 16:16:38 cram_md5_server authenticator failed for ([192.168.0.3]) [75.148.200.89]: 535 Incorrect authentication data (set_id=charles)"
	metric := NewAuthenticatorFailedMetric()
	if "535 Incorrect authentication data" != metric.parseErrorMessage(line) {
		t.Fail()
	}
}

func TestParseErrorMessageFailed(t *testing.T) {
	line := "2016-01-21 16:16:38 authenticator failed for ([192.168.0.3]) [75.148.200.89]"
	metric := NewAuthenticatorFailedMetric()
	if "unknown error message" != metric.parseErrorMessage(line) {
		t.Fail()
	}
}
