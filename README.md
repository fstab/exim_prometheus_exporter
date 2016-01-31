exim_prometheus_exporter
------------------------

[Prometheus](http://prometheus.io) exporter for monitoring the [Exim](http://www.exim.org/) Mail Transfer Agent (MTA).

Like [eximstats](http://www.exim.org/exim-html-current/doc/html/spec_html/ch-exim_utilities.html#SECTmailstat), the `exim_prometheus_exporter` generates statistics from Exim's `mainlog` file. It does not interact with Exim directly.

Usage
-----

```bash
exim_prometheus_exporter -mainlog /path/to/mainlog
```

By default, metrics are provided on [https://localhost:8443/metrics](https://localhost:8443/metrics).
Type `exim_prometheus_exporter -h` to see a list of command line options.


In order to access _https_ metris from the Prometheus server, add something like this to the server's `prometheus.yml` configuration:

```yaml
  - job_name: "exim-server"
    scrape_interval: "60s"
    scheme: "https"
    tls_config:
      insecure_skip_verify: true
    target_groups:
    - targets: ['mail.example.com:8443']
```

Installation
------------

Make sure [Go](https://golang.org) is installed and the `GOPATH` environment variable is set, then run

```bash
go get github.com/fstab/exim_prometheus_exporter
```

The executable will be created in `$GOPATH/bin`.

Status
------

This is just a proof-of-concept. Only `rejected RCPT` messages are recognized, all other messages are not implemented yet.
If this turns out to be useful, let's build on this example and add more messages.
