package fasthttp

import (
	"context"
	"fmt"
	"github.com/SkyAPM/go2sky"
	"github.com/valyala/fasthttp"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	"strconv"
	"time"
)

const (
	componentIDGoHttpClient = 5005
)

type ClientSpanHelper struct {
	tracer *go2sky.Tracer
}

func NewClientSpanHelper(trace *go2sky.Tracer) *ClientSpanHelper {
	return &ClientSpanHelper{
		tracer: trace,
	}
}

func (fc *ClientSpanHelper) CreateSpan(ctx context.Context, req *fasthttp.Request) (span go2sky.Span, err error) {
	if fc.tracer == nil {
		return
	}
	span, err = fc.tracer.CreateExitSpan(ctx, getOperationNameForReq(req), string(req.Host()), func(headerKey, headerValue string) error {
		reqHeader := &req.Header
		reqHeader.Set(headerKey, headerValue)
		return nil
	})
	if err != nil {
		return
	}
	span.SetComponent(componentIDGoHttpClient)
	span.Tag(go2sky.TagHTTPMethod, string(req.Header.Method()))
	span.Tag(go2sky.TagURL, req.URI().String())
	span.SetSpanLayer(agentv3.SpanLayer_Http)
	return
}

func (fc *ClientSpanHelper) EndSpan(span go2sky.Span, res *fasthttp.Response) {
	if span == nil {
		return
	}
	if res.StatusCode() >= 400 {
		span.Error(time.Now(), string(res.Body()))
	}
	span.Tag(go2sky.TagStatusCode, strconv.Itoa(res.StatusCode()))
	span.End()
}

func (fc *ClientSpanHelper) EndSpanError(span go2sky.Span, err error) {
	if span == nil {
		return
	}
	now := time.Now()
	span.Error(now, err.Error())
	span.End()
}

func getOperationNameForReq(req *fasthttp.Request) string {
	return fmt.Sprintf("%s:%s", req.Header.Method(), req.URI().Path())
}
