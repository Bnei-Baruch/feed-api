package instrumentation

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var Stats = new(Collectors)

type Collectors struct {
	RequestDurationHistogram    *prometheus.HistogramVec
	RecommendCounter            *prometheus.CounterVec
	RecommendSelectedCounter    *prometheus.CounterVec
	SearchCounter               prometheus.Counter
	SearchSelectedCounter       prometheus.Counter
	AutocompleteCounter         prometheus.Counter
	AutocompleteSelectedCounter prometheus.Counter
	SearchSelectedRankHistogram prometheus.Histogram
	EntriesCounterVec           *prometheus.CounterVec
}

func (c *Collectors) Init() {
	c.RequestDurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "feed_api",
		Subsystem: "api",
		Name:      "request_duration",
		Help:      "Time (in seconds) spent serving HTTP requests.",
	}, []string{"method", "route", "status_code"})

	c.RecommendCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "feed_api",
		Subsystem: "recommend",
		Name:      "recommend",
		Help:      "Counts number of recommendation served to users.",
	}, []string{"ab"})

	c.RecommendSelectedCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "feed_api",
		Subsystem: "recommend",
		Name:      "recommend_selected",
		Help:      "Counts number of recommendation selections by users.",
	}, []string{"ab"})

	c.SearchCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "feed_api",
		Subsystem: "search",
		Name:      "search",
		Help:      "Counts number of searches served to users.",
	})

	c.SearchSelectedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "feed_api",
		Subsystem: "search",
		Name:      "search_selected",
		Help:      "Counts number of search results selected by users.",
	})

	c.AutocompleteCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "feed_api",
		Subsystem: "search",
		Name:      "autocomplete",
		Help:      "Counts number of autocompletes served to users.",
	})

	c.AutocompleteSelectedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "feed_api",
		Subsystem: "search",
		Name:      "autocomplete_selected",
		Help:      "Counts number of autocomplete suggestions selected by users.",
	})

	c.SearchSelectedRankHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "feed_api",
		Subsystem: "search",
		Name:      "search_selected_rank",
		Help:      "Histogram of search selected rank.",
		Buckets: []float64{
			float64(0),
			float64(1),
			float64(2),
			float64(3),
			float64(4),
			float64(5),
			float64(6),
			float64(7),
			float64(8),
			float64(9),
			float64(10),
			float64(11),
			float64(12),
			float64(13),
			float64(14),
			float64(15),
			float64(16),
			float64(17),
			float64(18),
			float64(19),
		},
	})

	c.EntriesCounterVec = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "feed_api",
		Subsystem: "chronicles",
		Name:      "entries",
		Help:      "Counts number of entries.",
	}, []string{"event_type"})

	prometheus.MustRegister(c.RequestDurationHistogram)
	prometheus.MustRegister(c.RecommendCounter)
	prometheus.MustRegister(c.RecommendSelectedCounter)
	prometheus.MustRegister(c.SearchCounter)
	prometheus.MustRegister(c.SearchSelectedCounter)
	prometheus.MustRegister(c.AutocompleteCounter)
	prometheus.MustRegister(c.AutocompleteSelectedCounter)
	prometheus.MustRegister(c.SearchSelectedRankHistogram)
	prometheus.MustRegister(c.EntriesCounterVec)

	prometheus.MustRegister(collectors.NewBuildInfoCollector())
}

func (c *Collectors) Reset() {
	c.RequestDurationHistogram.Reset()
}
