package scrapesexp

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "nacos"
	subSystem = "service_exporter"
)

type Exporter struct {
	scrapes prometheus.Counter
}

func New() *Exporter {
	return &Exporter{
		scrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(namespace, subSystem, "total_scrapes"),
			Help: "nacos service total scrapes.",
		}),
	}
}
