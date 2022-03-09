package redis_hook

import (
	"context"
	"github.com/SkyAPM/go2sky"
	rv7 "github.com/go-redis/redis/v7"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	"strconv"
	"strings"
	"time"
)

// https://raw.githubusercontent.com/apache/skywalking/master/oap-server/server-starter/src/main/resources/component-libraries.yml
const componentID = 7

type tracerHook struct {
	tracer *go2sky.Tracer
	opts   *options
}

type Option func(*options)

type options struct {
	peer       string
	Db         string
	logError   bool
	reportArgs bool
}

func WithAddr(addr string) Option {
	return func(o *options) {
		o.peer = addr
	}
}

func WithDb(db int) Option {
	return func(o *options) {
		o.Db = strconv.Itoa(db)
	}
}

func WithLogError() Option {
	return func(o *options) {
		o.logError = true
	}
}

func WithArgsReport() Option {
	return func(o *options) {
		o.reportArgs = true
	}
}

func NewHook(tracer *go2sky.Tracer, opts ...Option) rv7.Hook {

	options := &options{}
	for _, o := range opts {
		o(options)
	}

	return &tracerHook{
		tracer: tracer,
		opts:   options,
	}
}

func noopInjector(headerKey, headerValue string) error {
	return nil
}

func (h *tracerHook) BeforeProcess(ctx context.Context, cmd rv7.Cmder) (context.Context, error) {
	span, nCtx, err := h.tracer.CreateExitSpanWithContext(ctx, cmd.Name(), h.opts.peer, noopInjector)

	if h.opts.reportArgs {
		span.Tag(go2sky.TagDBStatement, cmd.String())
	} else {
		span.Tag(go2sky.TagDBStatement, cmd.Name())
	}

	span.Tag(go2sky.TagDBInstance, h.opts.Db)
	span.Tag(go2sky.TagDBType, "Redis")
	span.SetComponent(componentID)
	span.SetSpanLayer(agentv3.SpanLayer_Cache)
	return nCtx, err
}

func (h *tracerHook) AfterProcess(ctx context.Context, cmd rv7.Cmder) error {
	span := go2sky.ActiveSpan(ctx)
	if span == nil {
		return nil
	}
	cmdErr := cmd.Err()
	if cmdErr != nil {
		now := time.Now()
		if h.opts.logError {
			span.Error(now, cmdErr.Error())
		} else {
			span.Error(now)
		}
	}
	span.End()
	return nil
}

func (h *tracerHook) BeforeProcessPipeline(ctx context.Context, cmds []rv7.Cmder) (context.Context, error) {

	opName := quickJoin(cmds, "|", func(obj rv7.Cmder) string {
		return obj.Name()
	})

	span, nCtx, err := h.tracer.CreateExitSpanWithContext(ctx, opName, h.opts.peer, noopInjector)

	if h.opts.reportArgs {
		stmt := quickJoin(cmds, "|", func(obj rv7.Cmder) string {
			return obj.String()
		})
		span.Tag(go2sky.TagDBStatement, stmt)
	} else {
		span.Tag(go2sky.TagDBStatement, opName)
	}
	span.Tag(go2sky.TagDBInstance, h.opts.Db)
	span.Tag(go2sky.TagDBType, "Redis")
	span.SetComponent(componentID)
	span.SetSpanLayer(agentv3.SpanLayer_Cache)
	return nCtx, err

}

type stringer func(obj rv7.Cmder) string

func quickJoin(objs []rv7.Cmder, sep string, stringer stringer) string {
	sb := &strings.Builder{}
	for i, obj := range objs {
		sb.WriteString(stringer(obj))
		if i < len(objs)-1 {
			sb.WriteString(sep)
		}
	}
	return sb.String()
}

func (h *tracerHook) AfterProcessPipeline(ctx context.Context, cmds []rv7.Cmder) error {
	span := go2sky.ActiveSpan(ctx)
	if span == nil {
		return nil
	}

	errs := quickJoinError(cmds, "|")
	if errs != "" {
		now := time.Now()
		if h.opts.logError {
			span.Error(now, errs)
		} else {
			span.Error(now)
		}
	}

	span.End()
	return nil
}

func quickJoinError(objs []rv7.Cmder, sep string) string {
	sb := &strings.Builder{}
	for _, obj := range objs {
		cmdErr := obj.Err()
		if cmdErr != nil {
			if sb.Len() > 0 {
				sb.WriteString(sep)
			}
			sb.WriteString(cmdErr.Error())
		}
	}
	return sb.String()
}
