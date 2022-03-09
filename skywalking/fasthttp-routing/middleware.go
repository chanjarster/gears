package fasthttp_routing

import (
	"context"
	"fmt"
	"github.com/SkyAPM/go2sky"
	routing "github.com/qiangxue/fasthttp-routing"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	"strconv"
	"time"
)

const (
	componentIDGoHttpServer = 5004
	componentIDGoHttpClient = 5005
	swCtxKey                = "_swCtxKey_"
)

var (
	GetSwContext = func(rCtx *routing.Context) context.Context {
		if swContext, ok := rCtx.Get(swCtxKey).(context.Context); ok && swContext != nil {
			return swContext
		}
		return context.Background()
	}
	GetOperationName = func(rCtx *routing.Context) string {
		return fmt.Sprintf("%s:%s", rCtx.Method(), rCtx.Path())
	}
	SetSwContext = func(rCtx *routing.Context, spanCtx context.Context) {
		rCtx.Set(swCtxKey, spanCtx)
	}
)

type options struct {
	logError   bool
	log4xxBody bool
}

type Option func(opts *options)

func WithLogError() Option {
	return func(opts *options) {
		opts.logError = true
	}
}

func WithLog4xxBody() Option {
	return func(opts *options) {
		opts.log4xxBody = true
	}
}

func NewMiddleware(tracer *go2sky.Tracer, opts ...Option) routing.Handler {
	if tracer == nil {
		return func(c *routing.Context) error {
			return c.Next()
		}
	}

	options := &options{}
	for _, o := range opts {
		o(options)
	}

	return func(rCtx *routing.Context) error {

		swContext := GetSwContext(rCtx)

		span, spanCtx, err := tracer.CreateEntrySpan(swContext, GetOperationName(rCtx), func(key string) (string, error) {
			return string(rCtx.Request.Header.Peek(key)), nil
		})
		if err != nil {
			return rCtx.Next()
		}
		span.SetComponent(componentIDGoHttpServer)
		span.Tag(go2sky.TagHTTPMethod, string(rCtx.RequestCtx.Method()))
		span.Tag(go2sky.TagURL, rCtx.RequestCtx.URI().String())
		span.SetSpanLayer(agentv3.SpanLayer_Http)

		SetSwContext(rCtx, spanCtx)

		err = rCtx.Next()
		if err != nil {
			if options.logError {
				span.Error(time.Now(), err.Error())
			} else {
				span.Error(time.Now())
			}
		} else if err = rCtx.Err(); err != nil {
			if options.logError {
				span.Error(time.Now(), err.Error())
			} else {
				span.Error(time.Now())
			}
		} else if rCtx.Response.StatusCode() >= 400 {
			if options.log4xxBody {
				span.Error(time.Now(), string(rCtx.Response.Body()))
			} else {
				span.Error(time.Now())
			}
		}
		span.Tag(go2sky.TagStatusCode, strconv.Itoa(rCtx.Response.StatusCode()))
		span.End()
		return nil
	}

}
