package ion

// TODO: move it from here when metric is implemented in kit

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
)

var (
	prometheusGaugeRooms = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ion_cluster_sessions",
			Help: "Number of currently active rooms on this node",
		},
	)

	prometheusGaugeClients = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ion_cluster_clients",
			Help: "Number of currently active websockets on this node",
		},
	)
	prometheusGaugeProxyClients = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ion_cluster_proxy_clients",
			Help: "Number of currently active proxied websockets on this node",
		},
	)
)

func init() {
	prometheus.MustRegister(prometheusGaugeRooms)
	prometheus.MustRegister(prometheusGaugeClients)
	prometheus.MustRegister(prometheusGaugeProxyClients)
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
}

func metricsHandler() http.Handler {
	return promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	)
}

// MetricsGetActiveClientsCount returns number of active clients connected to this node
func MetricsGetActiveClientsCount() int {
	clientCount := dto.Metric{}
	prometheusGaugeClients.Write(&clientCount)
	proxyCount := dto.Metric{}
	prometheusGaugeProxyClients.Write(&proxyCount)

	active := clientCount.GetGauge().GetValue() + proxyCount.GetGauge().GetValue()
	return int(active)
}
