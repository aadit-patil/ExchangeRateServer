package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	CacheHits = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_hits_total",
		Help: "Total number of cache hits",
	})
	CacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_misses_total",
		Help: "Total number of cache misses",
	})
	APIRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "external_api_requests_total",
		Help: "Total external API requests made",
	})
	DBQueries = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "db_queries_total",
		Help: "Total DB queries made",
	})
)

func Init() {
	prometheus.MustRegister(CacheHits, CacheMisses, APIRequests, DBQueries)
}
