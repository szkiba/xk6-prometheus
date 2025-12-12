// Package prometheus implements a Prometheus output for k6.
package prometheus

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/schema"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/szkiba/xk6-prometheus/internal"
	"go.k6.io/k6/output"
)

const defaultPort = 5656

// Register the extensions on module initialization.
func init() {
	output.RegisterExtension("prometheus", New)
}

type options struct {
	Port                int
	Host                string
	Subsystem           string
	Namespace           string
	UseHistogramForTime string
}

// Output is the Prometheus output implementation.
type Output struct {
	*internal.PrometheusAdapter

	addr   string
	arg    string
	logger logrus.FieldLogger
}

// New creates a new Prometheus output instance.
func New(params output.Params) (output.Output, error) { //nolint:ireturn
	registry, ok := prometheus.DefaultRegisterer.(*prometheus.Registry)
	if !ok {
		registry = prometheus.NewRegistry()
	}

	out := &Output{
		PrometheusAdapter: internal.NewPrometheusAdapter(registry, params.Logger, "", ""),
		arg:               params.ConfigArgument,
		logger:            params.Logger,
	}

	return out, nil
}

// Description implements output.Output.
func (o *Output) Description() string {
	return fmt.Sprintf("prometheus (%s)", o.addr)
}

func isOptTrue(s string) bool {
	return s == "true" || s == "yes"
}

// Start implements output.Output.
func (o *Output) Start() error {
	opts, err := getopts(o.arg)
	if err != nil {
		return err
	}

	o.Namespace = opts.Namespace
	o.Subsystem = opts.Subsystem
	o.addr = fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	o.UseHistogramForTime = isOptTrue(opts.UseHistogramForTime)

	listener, err := new(net.ListenConfig).Listen(context.TODO(), "tcp", o.addr)
	if err != nil {
		return err
	}

	go func() {
		server := &http.Server{Handler: o.Handler(), ReadHeaderTimeout: time.Second} //nolint:exhaustruct

		if err := server.Serve(listener); err != nil {
			o.logger.Error(err)
		}
	}()

	return nil
}

// Stop implements output.Output.
func (o *Output) Stop() error {
	return nil
}

func getopts(query string) (*options, error) {
	opts := &options{
		Port:                defaultPort,
		Host:                "",
		Namespace:           "",
		Subsystem:           "",
		UseHistogramForTime: "no",
	}

	if query == "" {
		return opts, nil
	}

	values, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}

	decoder := schema.NewDecoder()

	if err = decoder.Decode(opts, values); err != nil {
		return nil, err
	}

	return opts, nil
}
