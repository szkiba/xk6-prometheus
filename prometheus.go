package prometheus

import (
	"strings"

	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
	"github.com/loadimpact/k6/output"
	"github.com/loadimpact/k6/stats"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Register the extensions on module initialization.
func init() {
	output.RegisterExtension("prometheus", New)
}

type options struct {
	Port      int
	Host      string
	Subsystem string
	Namespace string
}

type outputImpl struct {
	params  *output.Params
	options *options
	metrics map[string]interface{}
}

func New(params output.Params) (output.Output, error) {
	o := &outputImpl{params: &params, metrics: make(map[string]interface{})}

	return o, nil
}

func (o *outputImpl) Description() string {
	return fmt.Sprintf("prometheus (%s:%d)", o.options.Host, o.options.Port)
}

func getopts(qs string) (*options, error) {
	opts := &options{
		Port:      5656,
		Host:      "",
		Namespace: "",
		Subsystem: "",
	}

	if qs == "" {
		return opts, nil
	}

	v, err := url.ParseQuery(qs)
	if err != nil {
		return nil, err
	}

	decoder := schema.NewDecoder()

	if err = decoder.Decode(opts, v); err != nil {
		return nil, err
	}

	return opts, nil
}

func (o *outputImpl) Start() (err error) {

	if o.options, err = getopts(o.params.ConfigArgument); err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", o.options.Host, o.options.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		if err := http.Serve(listener, promhttp.Handler()); err != nil {
			o.params.Logger.Error(err)
		}
	}()

	return nil
}

func (o *outputImpl) Stop() error {
	return nil
}

func (o *outputImpl) AddMetricSamples(samples []stats.SampleContainer) {
	for i := range samples {
		all := samples[i].GetSamples()
		for j := range all {
			o.handleSample(&all[j])
		}
	}
}

func (o *outputImpl) handleSample(sample *stats.Sample) {
	var handler func(*stats.Sample)

	switch sample.Metric.Type {
	case stats.Counter:
		handler = o.handleCounter
	case stats.Gauge:
		handler = o.handleGauge
	case stats.Rate:
		handler = o.handleRate
	case stats.Trend:
		handler = o.handleTrend
	default:
		o.params.Logger.Warnf("Unknown metric type: %v", sample.Metric.Type)
		return
	}

	handler(sample)
}

func (o *outputImpl) handleCounter(sample *stats.Sample) {
	if counter := o.getCounter(sample.Metric.Name, "k6 counter"); counter != nil {
		counter.Add(sample.Value)
	}
}

func (o *outputImpl) handleGauge(sample *stats.Sample) {
	if gauge := o.getGauge(sample.Metric.Name, "k6 gauge"); gauge != nil {
		gauge.Set(sample.Value)
	}
}

func (o *outputImpl) handleRate(sample *stats.Sample) {
	if histogram := o.getHistogram(sample.Metric.Name, "k6 rate", []float64{0}); histogram != nil {
		histogram.Observe(sample.Value)
	}
}

func (o *outputImpl) handleTrend(sample *stats.Sample) {
	if summary := o.getSummary(sample.Metric.Name, "k6 trend"); summary != nil {
		summary.Observe(sample.Value)
	}

	if gauge := o.getGauge(sample.Metric.Name+"_value", "k6 trend value"); gauge != nil {
		gauge.Set(sample.Value)
	}
}

func (o *outputImpl) getCounter(name string, helpSuffix string) (counter prometheus.Counter) {
	if col, ok := o.metrics[name]; ok {
		if c, tok := col.(prometheus.Counter); tok {
			counter = c
		}
	}

	if counter == nil {
		counter = prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: o.options.Namespace,
			Subsystem: o.options.Subsystem,
			Name:      name,
			Help:      helpFor(name, helpSuffix),
		})

		if err := prometheus.Register(counter); err != nil {
			o.params.Logger.Error(err)
			return nil
		}

		o.metrics[name] = counter
	}

	return counter
}

func (o *outputImpl) getGauge(name string, helpSuffix string) (gauge prometheus.Gauge) {
	if gau, ok := o.metrics[name]; ok {
		if g, tok := gau.(prometheus.Gauge); tok {
			gauge = g
		}
	}

	if gauge == nil {
		gauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: o.options.Namespace,
			Subsystem: o.options.Subsystem,
			Name:      name,
			Help:      helpFor(name, helpSuffix),
		})

		if err := prometheus.Register(gauge); err != nil {
			o.params.Logger.Error(err)
			return nil
		}

		o.metrics[name] = gauge
	}

	return gauge
}

func (o *outputImpl) getSummary(name string, helpSuffix string) (summary prometheus.Summary) {
	if sum, ok := o.metrics[name]; ok {
		if s, tok := sum.(prometheus.Summary); tok {
			summary = s
		}
	}

	if summary == nil {
		summary = prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace:  o.options.Namespace,
			Subsystem:  o.options.Subsystem,
			Name:       name,
			Help:       helpFor(name, helpSuffix),
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.001},
		})

		if err := prometheus.Register(summary); err != nil {
			o.params.Logger.Error(err)
			return nil
		}

		o.metrics[name] = summary
	}

	return summary
}

func (o *outputImpl) getHistogram(name string, helpSuffix string, buckets []float64) (histogram prometheus.Histogram) {
	if his, ok := o.metrics[name]; ok {
		if h, tok := his.(prometheus.Histogram); tok {
			histogram = h
		}
	}

	if histogram == nil {
		histogram = prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: o.options.Namespace,
			Subsystem: o.options.Subsystem,
			Name:      name,
			Help:      helpFor(name, helpSuffix),
			Buckets:   buckets,
		})

		if err := prometheus.Register(histogram); err != nil {
			o.params.Logger.Error(err)
			return nil
		}

		o.metrics[name] = histogram
	}

	return histogram
}

func helpFor(name string, helpSuffix string) string {
	if h, ok := builtinMetrics[name]; ok {
		return h
	}

	if h, ok := builtinMetrics[strings.TrimSuffix(name, "_value")]; ok {
		return h + " (value)"
	}

	return name + " " + helpSuffix
}

var (
	builtinMetrics = map[string]string{
		"vus":                "Current number of active virtual users",
		"vus_max":            "Max possible number of virtual users",
		"iterations":         "The aggregate number of times the VUs in the test have executed",
		"iteration_duration": "The time it took to complete one full iteration",
		"dropped_iterations": "The number of iterations that could not be started",
		"data_received":      "The amount of received data",
		"data_sent":          "The amount of data sent",
		"checks":             "The rate of successful checks",

		"http_reqs":                "How many HTTP requests has k6 generated, in total",
		"http_req_blocked":         "Time spent blocked  before initiating the request",
		"http_req_connecting":      "Time spent establishing TCP connection",
		"http_req_tls_handshaking": "Time spent handshaking TLS session",
		"http_req_sending":         "Time spent sending data",
		"http_req_waiting":         "Time spent waiting for response",
		"http_req_receiving":       "Time spent receiving response data",
		"http_req_duration":        "Total time for the request",
		"http_req_failed":          "The rate of failed requests",
	}
)
