package middleware

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

const (
	ServiceMainName    = "main"
	ServiceAuthName    = "auth"
	ServiceUserName    = "user"
	ServiceCreatorName = "creator"
)

var (
	UUIDRegExp = regexp.MustCompile(`[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)
)

const (
	ServiceName = "ServiceName"
	FullTime    = "Duration"
	URL         = "Url"
	Method      = "Method"
	StatusCode  = "StatusCode"
)

type writer struct {
	http.ResponseWriter
	statusCode int
}

func NewWriter(w http.ResponseWriter) *writer {
	return &writer{w, http.StatusOK}
}

func (w *writer) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

type MetricsMiddleware struct {
	metric    *prometheus.GaugeVec
	counter   *prometheus.CounterVec   //количество ошибок
	durations *prometheus.HistogramVec //сколько выполняются различные запросы
	errors    *prometheus.CounterVec
	name      string
}

func NewMetricsMiddleware() *MetricsMiddleware {
	return &MetricsMiddleware{}
}

func (m *MetricsMiddleware) ServerMetricsInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	start := time.Now()
	h, err := handler(ctx, req)
	tm := time.Since(start)

	m.metric.With(prometheus.Labels{
		URL:         "",
		ServiceName: m.name,
		StatusCode:  "OK",
		Method:      info.FullMethod,
		FullTime:    tm.String(),
	}).Inc()

	m.durations.With(prometheus.Labels{URL: info.FullMethod}).Observe(tm.Seconds())

	m.counter.With(prometheus.Labels{URL: info.FullMethod}).Inc()

	return h, err

}

func (m *MetricsMiddleware) Register(name string) {
	m.name = name
	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: fmt.Sprintf("SLO for service %s", name),
		},
		[]string{
			ServiceName, URL, Method, StatusCode, FullTime,
		})

	m.metric = gauge

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hits",
			Help: "Number of all requests.",
		}, []string{URL})
	m.counter = counter

	hist := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "durations_stats",
		Help:    "durations_stats",
		Buckets: prometheus.LinearBuckets(0, 1, 10),
	}, []string{URL})
	m.durations = hist

	errs := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "errors_hits",
		Help: "Number of all errors.",
	}, []string{URL})

	m.errors = errs
	rand.Seed(time.Now().Unix())
	prometheus.MustRegister(m.metric)
	prometheus.MustRegister(m.counter)
	prometheus.MustRegister(m.durations)
	prometheus.MustRegister(m.errors)
}

func (m *MetricsMiddleware) LogMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		start := time.Now()

		wrapper := NewWriter(w)

		next.ServeHTTP(wrapper, r.WithContext(ctx))

		tm := time.Since(start)

		bytesUrl := []byte(r.URL.Path)
		urlWithCuttedUUID := UUIDRegExp.ReplaceAll(bytesUrl, []byte("<uuid>"))

		m.metric.With(prometheus.Labels{
			ServiceName: m.name,
			URL:         string(urlWithCuttedUUID),
			Method:      r.Method,
			StatusCode:  fmt.Sprintf("%d", wrapper.statusCode),
			FullTime:    tm.String(),
		}).Inc()

		m.durations.With(prometheus.Labels{URL: string(urlWithCuttedUUID)}).Observe(tm.Seconds())

		if wrapper.statusCode != http.StatusOK {
			m.errors.With(prometheus.Labels{URL: string(urlWithCuttedUUID)}).Inc()
		}
		m.counter.With(prometheus.Labels{URL: string(urlWithCuttedUUID)}).Inc()
	})
}
