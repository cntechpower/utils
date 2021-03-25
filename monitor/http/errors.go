package http

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func SetSpanCode(span opentracing.Span, code int) {
	span.SetTag("response.code", code)
	if code >= 400 {
		ext.Error.Set(span, true)
	}
}
