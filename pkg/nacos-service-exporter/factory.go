package nacos_service_exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"nacos-service-discovery-controller/pkg/logger"
	"nacos-service-discovery-controller/pkg/nacos"
)

const (
	metricsNamespace = "nacos"
	metricsSubSystem = "service_discovery"
)

const (
	metricsNamespaceScrapesStatus = iota
	metricsNamespaceCount
	metricsServiceScrapesStatus
	metricsServiceCount
	metricsInstanceCount
)

var (
	namespaceScrapesStatusLabels []string
	namespaceCountLabels         []string

	serviceScrapesStatusLabels []string
	serviceCountLabels         = []string{
		"namespace_id", "namespace_show_name", "namespace_description", "details",
	}

	instanceCountLabels = []string{
		"namespace_id", "namespace_show_name", "namespace_description", "service_name", "details",
	}
)

type StatisticsMetrics struct {
	Type  prometheus.ValueType
	Desc  *prometheus.Desc
	Value func(m float64) float64
}

type Exporter struct {
	*Collector

	statisticsMetrics map[int]StatisticsMetrics
}

func New(config *nacos.Config) *Exporter {
	c, err := NewCollector(config)
	if err != nil {
		logger.Fatal("create collector failed", zap.Error(err))
	}
	if err := c.Run(); err != nil {
		logger.Fatal("run collector failed", zap.Error(err))
	}

	m := make(map[int]StatisticsMetrics, 5)
	// nacos_service_discovery_namespace_scrapes_status
	m[metricsNamespaceScrapesStatus] = StatisticsMetrics{
		Type: prometheus.GaugeValue,
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(metricsNamespace, metricsSubSystem, "namespace_scrapes_status"),
			"nacos service discovery namespace scrapes status.",
			namespaceScrapesStatusLabels, nil,
		),
		Value: func(m float64) float64 {
			return m
		},
	}
	// nacos_service_discovery_namespace_count
	m[metricsNamespaceCount] = StatisticsMetrics{
		Type: prometheus.GaugeValue,
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(metricsNamespace, metricsSubSystem, "namespace_count"),
			"nacos service discovery namespace count.",
			namespaceCountLabels, nil,
		),
		Value: func(m float64) float64 {
			return m
		},
	}
	// nacos_service_discovery_service_scrapes_status
	m[metricsServiceScrapesStatus] = StatisticsMetrics{
		Type: prometheus.GaugeValue,
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(metricsNamespace, metricsSubSystem, "service_scrapes_status"),
			"nacos service discovery service scrapes status.",
			serviceScrapesStatusLabels, nil,
		),
		Value: func(m float64) float64 {
			return m
		},
	}
	// nacos_service_discovery_service_count
	m[metricsServiceCount] = StatisticsMetrics{
		Type: prometheus.GaugeValue,
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(metricsNamespace, metricsSubSystem, "service_count"),
			"nacos service discovery service count.",
			serviceCountLabels, nil,
		),
		Value: func(m float64) float64 {
			return m
		},
	}
	// nacos_service_discovery_instance_count
	m[metricsInstanceCount] = StatisticsMetrics{
		Type: prometheus.GaugeValue,
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(metricsNamespace, metricsSubSystem, "instance_count"),
			"nacos service discovery instance count.",
			instanceCountLabels, nil,
		),
		Value: func(m float64) float64 {
			return m
		},
	}

	return &Exporter{
		Collector:         c,
		statisticsMetrics: m,
	}
}
