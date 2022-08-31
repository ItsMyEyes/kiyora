package tracing

import (
	"context"
	"fmt"
	"io"

	"myself_framwork/library/log"

	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
)

// LogrusAdapter - an adapter to log span info
type LogrusAdapter struct {
	InfoLevel bool
}

// Error - logrus adapter for span errors
func (l LogrusAdapter) Error(msg string) {
	logrus.Error(msg)
}

// Infof - logrus adapter for span info logging
func (l LogrusAdapter) Infof(msg string, args ...interface{}) {
	if l.InfoLevel {
		logrus.Infof(msg, args...)
	}
}

// Option - define options for NewJWTCache()
type Option func(*options)
type options struct {
	sampleProbability float64
	enableInfoLog     bool
}

// defaultOptions - some defs options to NewJWTCache()
var defaultOptions = options{
	sampleProbability: 0.0,
	enableInfoLog:     false,
}

// WithSampleProbability - optional sample probability
func WithSampleProbability(sampleProbability float64) Option {
	return func(o *options) {
		o.sampleProbability = sampleProbability
	}
}

// WithEnableInfoLog - optional: enable Info logging for tracing
func WithEnableInfoLog(enable bool) Option {
	return func(o *options) {
		o.enableInfoLog = enable
	}
}

type jaegerLoggerAdapter struct {
	logger log.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}

// InitTracing - init opentracing with options (WithSampleProbability, WithEnableInfoLog) defaults: constant sampling, no info logging
func NewInitTracing(serviceName string,
	metricsFactory metrics.Factory, logger log.Factory) opentracing.Tracer {
	cfg, err := config.FromEnv()
	if err != nil {
		logger.Bg().Fatal("cannot parse Jaeger env vars", zap.Error(err))
	}
	cfg.ServiceName = serviceName
	cfg.Sampler.Type = "const"
	cfg.Sampler.Param = 1

	jaegerLogger := jaegerLoggerAdapter{logger.Bg()}

	metricsFactory = metricsFactory.Namespace(metrics.NSOptions{Name: serviceName, Tags: nil})
	tracer, _, err := cfg.NewTracer(
		config.Logger(jaegerLogger),
		config.Metrics(metricsFactory),
		config.Observer(rpcmetrics.NewObserver(metricsFactory, rpcmetrics.DefaultNameNormalizer)),
	)
	if err != nil {
		logger.Bg().Fatal("cannot initialize Jaeger Tracer", zap.Error(err))
	}
	return tracer

}

// InitTracing - init opentracing with options (WithSampleProbability, WithEnableInfoLog) defaults: constant sampling, no info logging
func InitTracing(serviceName string, tracingAgentHostPort string, opt ...Option) (
	tracer opentracing.Tracer,
	reporter jaeger.Reporter,
	closer io.Closer,
	err error) {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}
	// factory := jaegerprom.New()
	// metrics := jaeger.NewMetrics(factory, map[string]string{"lib": "jaeger"})

	transport, err := jaeger.NewUDPTransport(tracingAgentHostPort, 0)
	if err != nil {
		return tracer, reporter, closer, err
	}

	logAdapt := LogrusAdapter{InfoLevel: opts.enableInfoLog}
	reporter = jaeger.NewCompositeReporter(
		jaeger.NewLoggingReporter(logAdapt),
		jaeger.NewRemoteReporter(transport,
			// jaeger.ReporterOptions.Metrics(metrics),
			jaeger.ReporterOptions.Logger(logAdapt),
		),
	)
	sampler := jaeger.NewConstSampler(true)
	//if opts.sampleProbability > 0 {
	//	fmt.Println("probable")
	//	sampler, err = jaeger.NewProbabilisticSampler(opts.sampleProbability)
	//}

	tracer, closer = jaeger.NewTracer(serviceName,
		sampler,
		reporter,
		// jaeger.TracerOptions.Metrics(metrics),
	)
	fmt.Println("Init Tracing Success ")
	return tracer, reporter, closer, nil
}

func Introduce_span(ctx context.Context, spanName string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, spanName)
	return span, ctx
}
