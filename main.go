package main

import (
	"flag"
	"fmt"
	"github.com/ActiveState/tail"
	"github.com/fstab/exim_prometheus_exporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"net/http"
	"os"
)

// cert and key created with openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -nodes

const CERT = `-----BEGIN CERTIFICATE-----
MIIDtTCCAp2gAwIBAgIJAP9eE4ZtnJZnMA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTYwMTMxMTUzNzAxWhcNMTYwMzAxMTUzNzAxWjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8GA1UEChMYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAuGMNuDwIdCbSNnOd0Bl3LejAgYjcXb6I2cdZCISt45CERO/mvn59MLYj
P7awhZuOBYhu00WshNfonbyZ7mNgTIIe8MuRIbHvqVhb2i8CwleorJhzT6cnfFf5
xjXj8DclSZDJAqQzthvGra7F67G38bKjl/0tx4T7Z4shp4d+M9to4zQp5x3xZ6hj
/3J9oZiMbAy8s+kODqIPHCsVjiCQqr/649tF5Fiq+UGRcOzR2471xKqhB37nMfAz
uoE5P6HENGN4K+fG8yJ7biBz063GZbcjopIj7RSZK9eZGfzGZm0NoqOUyjsyCKpS
0teK4Frw6Um8wTPkRdOypNvVgchDTQIDAQABo4GnMIGkMB0GA1UdDgQWBBRm4vRi
X6QudIwhw76HQRQQLj7knTB1BgNVHSMEbjBsgBRm4vRiX6QudIwhw76HQRQQLj7k
naFJpEcwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgTClNvbWUtU3RhdGUxITAfBgNV
BAoTGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZIIJAP9eE4ZtnJZnMAwGA1UdEwQF
MAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAK8lJOLaSgr0rGt3tLQ5lDG28L1tK2ki
TZNGXvbZOupdnQPC90gP0EiiNXDLIqqYYe12tA3YV6fQ+PhAaUdC413njHQGtbmn
2P/uLfHycrsttUpqmbWDF+uCdW/z42jztxMg4Ett2mMgt0LyQkyBxOCCi3Ia8cy9
GttqN9uvNsyfCVzBMC8DLpmnhBMtGJ2to0C/ktyzM8Z2t8I6F/RJUKiC/qAaqvHa
mrYEF+Y3rJCXlZoRCyd8j2lT2VSRFoOv+LgOV4aUgLd8Hw+epRJ/I+gO5YoQ1n6x
tEOmZmjlXW1eoFDTp2jnVN4gL7u4/d4B7F0kCltb3/3ZtWTv8AeTDSE=
-----END CERTIFICATE-----
`

const KEY = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAuGMNuDwIdCbSNnOd0Bl3LejAgYjcXb6I2cdZCISt45CERO/m
vn59MLYjP7awhZuOBYhu00WshNfonbyZ7mNgTIIe8MuRIbHvqVhb2i8CwleorJhz
T6cnfFf5xjXj8DclSZDJAqQzthvGra7F67G38bKjl/0tx4T7Z4shp4d+M9to4zQp
5x3xZ6hj/3J9oZiMbAy8s+kODqIPHCsVjiCQqr/649tF5Fiq+UGRcOzR2471xKqh
B37nMfAzuoE5P6HENGN4K+fG8yJ7biBz063GZbcjopIj7RSZK9eZGfzGZm0NoqOU
yjsyCKpS0teK4Frw6Um8wTPkRdOypNvVgchDTQIDAQABAoIBAQCGsf18v4Yha5aW
poD7Ww7/3455UgRBCwYnqQO2QE5S9ehZ/7JdKEPFyNgZHBj5kTf/fLoQ5k3vwVWx
nOwKBFh9q3R0zRCZP8XmvKBk04C9fZG/e6KI5n/mytGw5P89JNu9UOI2ZsNL3iCW
Eh2NXwcTrj7psc62eMO60R1lp4oe0ITyKEhIbHWK1cu+rRpmpIOOjdZK8fz296bX
HdM/tIw6OhXun6dEMRAY1pTy5Eznpvi0mIp+o8pGovm0KUokf87MTDi4GMnm14/7
znej5i6ETWYk1Yy1n/nzf+6OiAg66qu1TX76d/8w1JmY//LWpXuk3xEr/ZDBNMtk
/PrYswpBAoGBANm+iwfVssuyC5hcTva0InycU6jaKjPIt+Q7OWBAKR/EXj0CMpni
WYYR7M9Mcous2J1Rf8tTJr9K5RMheHY3zGn6qHYirUj8wDF/hlhvNG6RFME1U9NT
dfWBiD4JKDoTEzABDJQm6IroZpfN/89O1SuxzPuaLx4LVlY4eRP8ocTRAoGBANjI
MzPI+wlv2FFq9LjbNiD95lfIhp2xumrhmoplM5T9YlAiymtGU9fbA0GwKEqDlaqf
niXenkZhLclUbc6u85hOfmmoq5pEd6xk/OUstrd5dAxQMCRGlXMcEBYvazVbqUTp
5PbJ3FX7eKDxRi/AIKWV+DY2289yyF6BL0HjW+W9AoGAXlP8WNWL0lB8U3HRx3A7
7G2wlFqGo85VU6sQbRD+f8OK67UTBLUZAUqsoxVEHhwv7t8KlKOeCorAeCwsylHb
3SF4b00QcqkD/a14HsF2Hlv9eMHIYakrVcLaqb0/zwDKdCZQM7IzVVHed+8G3eER
2g75dRnTRZm1uj5WvYDY97ECgYEAuMC20pWhTYOiypDrDHjXAvsgywO9prwH8ntf
qD9j3MCufzmHZjHD1x1zAxLM4+SNM6Nht0ipf7XmvcVU6Gc2eEG9fvMffRSJIcXX
usGG34uFGdFlliUJzdbG5wF2zzzVYEQuvR2AyU7OmevHM3780+KibiIG6CAdIF3d
FrxcX8kCgYEAh+9YpiAVaYGDiAiqUMFRlQuDod/XN8HVreAHHn5G+A5O8yQMjB4T
/o2pN0JuECBcuVTcxl94w6kmfX1fQS3q7j+gZMzlO+eVweuLC19gO96RPUETxgeV
XLgD9hrDBrTbnKBHHQ6MHpT6ILi4w/e4+5XEUUOBf44ZJE71uRr4ZUA=
-----END RSA PRIVATE KEY-----
`

var (
	mainlog = flag.String("mainlog", "/var/log/exim4/mainlog", "Path to the exim4 mainlog file.")
	port    = flag.Int("port", 8443, "The port to listen on for HTTPS requests.")
	cert    = flag.String("cert", "", "Path to the SSL certificate file for the HTTPS server. (optional)")
	key     = flag.String("key", "", "Path to the SSL private key file for the HTTPS server. (optional)")
)

type Metric interface {
	Collector() prometheus.Collector
	Observe(logline string)
}

var (
	Metrics = []Metric{
		metrics.NewRejectedRcptMetric(),
		metrics.NewAuthenticatorFailedMetric(),
	}
)

func main() {
	flag.Parse()
	t := mustInitMainlogWatcher()
	for _, m := range Metrics {
		prometheus.MustRegister(m.Collector())
	}
	go runMainlogWatcher(t)
	runHttpServer()
}

func mustInitMainlogWatcher() *tail.Tail {
	t, err := tail.TailFile(*mainlog, tail.Config{
		MustExist: true,                  // Fail early if the file does not exist
		ReOpen:    true,                  // Reopen recreated files (tail -F)
		Follow:    true,                  // Continue looking for new lines (tail -f)
		Logger:    tail.DiscardingLogger, // Disable logging
		Location:  &tail.SeekInfo{0, 2},  // Start at the end of the file
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}
	return t
}

func runMainlogWatcher(t *tail.Tail) {
	for line := range t.Lines {
		if line.Err != nil {
			fmt.Errorf("Failed to tail %v: %v", *mainlog, line.Err.Error())
			os.Exit(-1)
		}
		for _, m := range Metrics {
			m.Observe(line.Text)
		}
	}
}

func runHttpServer() {
	if len(*cert) == 0 && len(*key) > 0 || len(*cert) > 0 && len(*key) == 0 {
		fmt.Fprintln(os.Stderr, "Syntax error: -cert and -key cannot be used without each other.")
		os.Exit(-1)
	}
	if len(*cert) == 0 || len(*key) == 0 {
		*cert = mustCreateTempFile("cert", []byte(CERT))
		*key = mustCreateTempFile("key", []byte(KEY))
		defer os.Remove(*cert)
		defer os.Remove(*key)
	}
	http.Handle("/metrics", prometheus.Handler())
	fmt.Printf("Starting server on https://localhost:%v/metrics\n", *port)
	err := http.ListenAndServeTLS(fmt.Sprintf(":%v", *port), *cert, *key, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}
}

func mustCreateTempFile(prefix string, data []byte) string {
	tempFile, err := ioutil.TempFile(os.TempDir(), prefix)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create temporary file: %v", err.Error())
		os.Exit(-1)
	}
	_, err = tempFile.Write(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write temporary file: %v", err.Error())
		os.Exit(-1)
	}
	err = tempFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close temporary file: %v", err.Error())
		os.Exit(-1)
	}
	return tempFile.Name()
}
