package skywalking

import (
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"time"
)

type noopReporter int

func (s noopReporter) Boot(service string, serviceInstance string, cdsWatchers []go2sky.AgentConfigChangeWatcher) {
}

func (s noopReporter) Send(spans []go2sky.ReportedSpan) {
}

func (s noopReporter) Close() {
}

func NoopReporter() go2sky.Reporter {
	return noopReporter(1)
}

func NewGrpcOrNoopReporter(c *Conf) (go2sky.Reporter, error) {
	if c == nil || c.ServerAddr == "" {
		return NoopReporter(), nil
	}

	reporter, err := reporter.NewGRPCReporter(c.ServerAddr,
		reporter.WithCDS(time.Minute),
		reporter.WithCheckInterval(time.Second*10),
	)
	if err != nil {
		return nil, err
	}

	return reporter, nil
}

func NewTracer(c *Conf, rep go2sky.Reporter) (*go2sky.Tracer, error) {
	return go2sky.NewTracer(c.ServiceNamespace+"::"+c.ServiceName,
		go2sky.WithReporter(rep),
		go2sky.WithInstance(c.ServiceInstance),
	)
}
