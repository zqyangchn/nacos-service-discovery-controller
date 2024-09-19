package scrapesexp

import (
	"github.com/prometheus/client_golang/prometheus"
)

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrapes.Inc()
	ch <- e.scrapes
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.scrapes.Desc()
}
