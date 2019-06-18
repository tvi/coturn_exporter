package coturn

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Active Sessions Metrics
	activeSessions = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "coturn_active_sessions_total",
		Help: "The total number of active sessions",
	})

	activeSessionsAge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "coturn_active_sessions_age_seconds",
		Help: "Average age of active sessions",
	})

	activeSessionsExpiration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "coturn_active_sessions_expiration_seconds",
		Help: "Average expiration of active sessions",
	})

	activeSessionsSentBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "coturn_active_sessions_transmit_bytes",
		Help: "Average sent bytes of active sessions",
	})

	activeSessionsRecvBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "coturn_active_sessions_receive_bytes",
		Help: "Average received bytes of active sessions",
	})

	activeSessionsTotalBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "coturn_active_sessions_total_bytes",
		Help: "Average send and received bytes of active sessions",
	})

	sessionBytesBuckets = []float64{
		1000,       // Metric KB
		1000000,    // Metric MB
		5000000,    // 5 MB
		10000000,   // 10 MB
		20000000,   // 20 MB
		50000000,   // 50 MB
		70000000,   // 70 MB
		100000000,  // 100 MB
		200000000,  // 200 MB
		500000000,  // 500 MB
		1000000000, // Metric GB
		10000000000,
	}

	sessionSecondsBuckets = []float64{
		1, // 1 second
		5,
		15,
		30,
		60,              // 1 minute
		2 * 60,          // 2 minutes
		3 * 60,          // 3 minutes
		4 * 60,          // 4 minutes
		5 * 60,          // 5 minutes
		7 * 60,          // 7 minutes
		10 * 60,         // 10 minutes
		15 * 60,         // 15 minutes
		20 * 60,         // 20 minutes
		30 * 60,         // 30 minutes
		45 * 60,         // 45 minutes
		1 * 60 * 60,     // 1 hour
		1*60*60 + 30*60, // 1.5 hour
		2 * 60 * 60,     // 2 hours
		2*60*60 + 30*60, // 2.5 hour
		3 * 60 * 60,     // 3 hours
		6 * 60 * 60,     // 6 hours
		12 * 60 * 60,    // 12 hours
		24 * 60 * 60,    // 24 hours
	}

	// All Sessions Metrics
	sessionsAge = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "coturn_sessions_age_seconds",
		Help:    "Histogram age of all finished sessions",
		Buckets: sessionSecondsBuckets,
	})

	sessionsSentBytes = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "coturn_sessions_transmit_bytes",
		Help:    "Sent bytes of all finished sessions",
		Buckets: sessionBytesBuckets,
	})

	sessionsRecvBytes = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "coturn_sessions_receive_bytes",
		Help:    "Received bytes of all finished sessions",
		Buckets: sessionBytesBuckets,
	})

	sessionsTotalBytes = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "coturn_sessions_total_bytes",
		Help:    "Total send and received bytes of all finished sessions",
		Buckets: sessionBytesBuckets,
	})
)

func init() {
	// Active Sessions
	prometheus.MustRegister(activeSessions)
	prometheus.MustRegister(activeSessionsAge)
	prometheus.MustRegister(activeSessionsExpiration)
	prometheus.MustRegister(activeSessionsSentBytes)
	prometheus.MustRegister(activeSessionsRecvBytes)
	prometheus.MustRegister(activeSessionsTotalBytes)

	// All Sesssions
	prometheus.MustRegister(sessionsAge)
	prometheus.MustRegister(sessionsSentBytes)
	prometheus.MustRegister(sessionsRecvBytes)
	prometheus.MustRegister(sessionsTotalBytes)

}
