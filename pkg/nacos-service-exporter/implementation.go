package nacos_service_exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"nacos-service-discovery-controller/pkg/logger"
)

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// nacos_service_discovery_namespace_scrapes_status
	if m, ok := e.statisticsMetrics[metricsNamespaceScrapesStatus]; ok {
		status := 0
		if e.namespaceScrapesSuccess {
			status = 1
		}
		ch <- prometheus.MustNewConstMetric(
			m.Desc,
			m.Type,
			m.Value(float64(status)),
		)
	}
	// nacos_service_discovery_namespace_count
	if m, ok := e.statisticsMetrics[metricsNamespaceCount]; ok {
		namespaces, err := e.GetNamespaces()
		if err != nil {
			logger.Warn("Failed to get namespaces", zap.Error(err))
			goto METRICSSERVICESCRAPESSTATUS // 不喜欢 else, 故 goto
		}
		ch <- prometheus.MustNewConstMetric(
			m.Desc,
			m.Type,
			m.Value(float64(len(namespaces))),
		)
	}

METRICSSERVICESCRAPESSTATUS:
	// nacos_service_discovery_service_scrapes_status
	if m, ok := e.statisticsMetrics[metricsServiceScrapesStatus]; ok {
		status := 0
		if e.serviceScrapesSuccess {
			status = 1
		}
		ch <- prometheus.MustNewConstMetric(
			m.Desc,
			m.Type,
			m.Value(float64(status)),
		)
	}
	// nacos_service_discovery_service_count
	if m, ok := e.statisticsMetrics[metricsServiceCount]; ok {
		services, err := e.GetServices()
		if err != nil {
			logger.Warn("Failed to get services", zap.Error(err))
			goto METRICSINSTANCECOUNT
		}
		for _, s := range services {
			details := ""
			if s.Error != nil {
				details = s.Error.Error()
			}
			ch <- prometheus.MustNewConstMetric(
				m.Desc,
				m.Type,
				m.Value(float64(len(s.services))),
				s.Id,
				s.ShowName,
				s.Description,
				details,
			)
		}
	}

METRICSINSTANCECOUNT:
	// nacos_service_discovery_instance_count
	if m, ok := e.statisticsMetrics[metricsInstanceCount]; ok {
		instancesCount, err := e.GetInstancesCount()
		if err != nil {
			return
		}
		for _, ic := range instancesCount {
			details := ""
			if ic.Error != nil {
				details = ic.Error.Error()
			}
			ch <- prometheus.MustNewConstMetric(
				m.Desc,
				m.Type,
				m.Value(float64(ic.Count)),
				ic.Id,
				ic.ShowName,
				ic.Description,
				ic.Name,
				details,
			)
		}
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, v := range e.statisticsMetrics {
		ch <- v.Desc
	}
}
