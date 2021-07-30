package instrumentation

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var Stats = new(Collectors)

type Collectors struct {
	RequestDurationHistogram *prometheus.HistogramVec
	RecommendCounter         prometheus.Counter
	RecommendSelectedCounter prometheus.Counter
}

func (c *Collectors) Init() {
	c.RequestDurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "feed_api",
		Subsystem: "api",
		Name:      "request_duration",
		Help:      "Time (in seconds) spent serving HTTP requests.",
	}, []string{"method", "route", "status_code"})

	c.RecommendCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "feed_api",
		Subsystem: "recommend",
		Name:      "recommend",
		Help:      "Counts number of recommendation served to users.",
	})

	c.RecommendSelectedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "feed_api",
		Subsystem: "recommend",
		Name:      "recommend_selected",
		Help:      "Counts number of recommendation selections by users.",
	})

	prometheus.MustRegister(c.RequestDurationHistogram)
	prometheus.MustRegister(c.RecommendCounter)
	prometheus.MustRegister(c.RecommendSelectedCounter)
	prometheus.MustRegister(collectors.NewBuildInfoCollector())
}

func (c *Collectors) Reset() {
	c.RequestDurationHistogram.Reset()
}
