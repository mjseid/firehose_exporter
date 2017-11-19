package collectors

import (
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/mjseid/firehose_exporter/cfinstanceinfoapi"
	"github.com/mjseid/firehose_exporter/metrics"
)

type ContainerMetricsCollector struct {
	namespace              string
	environment            string
	metricsStore           *metrics.Store
	cpuPercentageMetric    *prometheus.GaugeVec
	memoryBytesMetric      *prometheus.GaugeVec
	diskBytesMetric        *prometheus.GaugeVec
	memoryBytesQuotaMetric *prometheus.GaugeVec
	diskBytesQuotaMetric   *prometheus.GaugeVec
	appinfo                map[string]cfinstanceinfoapi.AppInfo
	amutex		       *sync.RWMutex
}

func NewContainerMetricsCollector(
	namespace string,
	environment string,
	metricsStore *metrics.Store,
	appinfo map[string]cfinstanceinfoapi.AppInfo,
	amutex *sync.RWMutex,
) *ContainerMetricsCollector {
	cpuPercentageMetric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   namespace,
			Subsystem:   container_metrics_subsystem,
			Name:        "cpu_percentage",
			Help:        "Cloud Foundry Firehose container metric: CPU used, on a scale of 0 to 100.",
			ConstLabels: prometheus.Labels{"environment": environment},
		},
                []string{"bosh_job_ip", "application_id", "instance_index", "app_name", "space", "org"},
	)

	memoryBytesMetric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   namespace,
			Subsystem:   container_metrics_subsystem,
			Name:        "memory_bytes",
			Help:        "Cloud Foundry Firehose container metric: bytes of memory used.",
			ConstLabels: prometheus.Labels{"environment": environment},
		},
                []string{"bosh_job_ip", "application_id", "instance_index", "app_name", "space", "org"},
	)

	diskBytesMetric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   namespace,
			Subsystem:   container_metrics_subsystem,
			Name:        "disk_bytes",
			Help:        "Cloud Foundry Firehose container metric: bytes of disk used.",
			ConstLabels: prometheus.Labels{"environment": environment},
		},
                []string{"bosh_job_ip", "application_id", "instance_index", "app_name", "space", "org"},
	)

	memoryBytesQuotaMetric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   namespace,
			Subsystem:   container_metrics_subsystem,
			Name:        "memory_bytes_quota",
			Help:        "Cloud Foundry Firehose container metric: maximum bytes of memory allocated to container.",
			ConstLabels: prometheus.Labels{"environment": environment},
		},
                []string{"bosh_job_ip", "application_id", "instance_index", "app_name", "space", "org"},
	)

	diskBytesQuotaMetric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   namespace,
			Subsystem:   container_metrics_subsystem,
			Name:        "disk_bytes_quota",
			Help:        "Cloud Foundry Firehose container metric: maximum bytes of disk allocated to container.",
			ConstLabels: prometheus.Labels{"environment": environment},
		},
                []string{"bosh_job_ip", "application_id", "instance_index", "app_name", "space", "org"},
	)

	return &ContainerMetricsCollector{
		namespace:              namespace,
		environment:            environment,
		metricsStore:           metricsStore,
		cpuPercentageMetric:    cpuPercentageMetric,
		memoryBytesMetric:      memoryBytesMetric,
		diskBytesMetric:        diskBytesMetric,
		memoryBytesQuotaMetric: memoryBytesQuotaMetric,
		diskBytesQuotaMetric:   diskBytesQuotaMetric,
		appinfo:                appinfo,
		amutex: 		amutex,
	}
}

func (c ContainerMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	c.amutex.Lock()

	c.cpuPercentageMetric.Reset()
	c.memoryBytesMetric.Reset()
	c.diskBytesMetric.Reset()
	c.memoryBytesQuotaMetric.Reset()
	c.diskBytesQuotaMetric.Reset()

	for _, containerMetric := range c.metricsStore.GetContainerMetrics() {
		Name := c.appinfo[containerMetric.ApplicationId].Name
		Space := c.appinfo[containerMetric.ApplicationId].Space
		Org := c.appinfo[containerMetric.ApplicationId].Org

		c.cpuPercentageMetric.WithLabelValues(
			containerMetric.IP,
			containerMetric.ApplicationId,
			strconv.Itoa(int(containerMetric.InstanceIndex)),
			Name,
                        Space,
                        Org,
		).Set(containerMetric.CpuPercentage)

		c.memoryBytesMetric.WithLabelValues(
			containerMetric.IP,
			containerMetric.ApplicationId,
			strconv.Itoa(int(containerMetric.InstanceIndex)),
			Name,
                        Space,
                        Org,
		).Set(float64(containerMetric.MemoryBytes))
		
		c.diskBytesMetric.WithLabelValues(
			containerMetric.IP,
			containerMetric.ApplicationId,
			strconv.Itoa(int(containerMetric.InstanceIndex)),
			Name,
                        Space,
                        Org,
		).Set(float64(containerMetric.DiskBytes))

		c.memoryBytesQuotaMetric.WithLabelValues(
			containerMetric.IP,
			containerMetric.ApplicationId,
			strconv.Itoa(int(containerMetric.InstanceIndex)),
			Name,
                        Space,
                        Org,
		).Set(float64(containerMetric.MemoryBytesQuota))

		c.diskBytesQuotaMetric.WithLabelValues(
			containerMetric.IP,
			containerMetric.ApplicationId,
			strconv.Itoa(int(containerMetric.InstanceIndex)),
			Name,
                        Space,
                        Org,
		).Set(float64(containerMetric.DiskBytesQuota))
	}
	
	c.cpuPercentageMetric.Collect(ch)
	c.memoryBytesMetric.Collect(ch)
	c.diskBytesMetric.Collect(ch)
	c.memoryBytesQuotaMetric.Collect(ch)
	c.diskBytesQuotaMetric.Collect(ch)

	c.amutex.Unlock()
}

func (c ContainerMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	c.cpuPercentageMetric.Describe(ch)
	c.memoryBytesMetric.Describe(ch)
	c.diskBytesMetric.Describe(ch)
	c.memoryBytesQuotaMetric.Describe(ch)
	c.diskBytesQuotaMetric.Describe(ch)
}
