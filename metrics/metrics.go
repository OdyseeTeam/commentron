package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// WSConnections is a metric to show the number of active web socket connections being handled
	WSConnections = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "commentron",
		Subsystem: "websocket",
		Name:      "connections",
		Help:      "Number of active web socket connections",
	}, []string{"claim"})

	// UserLoadOverall Number of active users
	UserLoadOverall = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "commentron",
		Subsystem: "apis",
		Name:      "user_load",
		Help:      "Number of active users",
	})

	// UserLoadByAPI Number of active calls by api
	UserLoadByAPI = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "commentron",
		Subsystem: "apis",
		Name:      "api_load",
		Help:      "Number of active calls by api",
	}, []string{"path"})

	// Durations The durations of the individual api calls
	Durations = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "commentron",
		Subsystem: "apis",
		Name:      "duration",
		Help:      "The durations of the individual api calls",
	}, []string{"path"})
)
