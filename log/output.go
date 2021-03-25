package log

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cntechpower/utils/tracing"
)

func logOutput(ctx context.Context, skip int, h *Header, level Level, format string, a ...interface{}) {
	file := ""
	line := 0
	if h.reportFileLine {
		file, line = getCaller(skip)
	}
	for _, l := range loggers {
		switch l.typ {
		case OutputTypeText:
			logOutputText(ctx, l, file, line, h, level, format, a...)
		case OutputTypeJson:
			logOutputStructured(ctx, l, file, line, h, level, format, a...)
		}
	}
}

func logOutputText(ctx context.Context, l *loggerWithConfig, file string, line int, h *Header, level Level, format string, a ...interface{}) {
	var s string
	if h.reportFileLine {
		s = fmt.Sprintf("[%s] <%s> |%s|%s| (%s:%v) %s",
			time.Now().Format("2006-01-02 15:04:05.000"), level, h, h.fields.String(),
			file, line, fmt.Sprintf(format, a...))
	} else {
		s = fmt.Sprintf("[%s] <%s> |%s|%s| %s",
			time.Now().Format("2006-01-02 15:04:05.000"), level, h, h.fields.String(),
			fmt.Sprintf(format, a...))
	}
	select {
	case l.buffer <- s:
		return
	case <-time.After(time.Millisecond):
		fmt.Printf("[%v] drop log because log buffer is full %v\n", l.typ, s)
	}
}

func logOutputStructured(ctx context.Context, l *loggerWithConfig, file string, line int, h *Header, level Level, format string, a ...interface{}) {
	nf := h.fields.DeepCopy()
	for k, v := range defaultFields {
		nf[k] = v
	}
	nf[fieldNameTime] = time.Now().Format(time.RFC3339)
	if h.reportFileLine {
		nf[fieldNameFileName] = file
		nf[fieldNameFileLine] = line
	}
	nf[fieldNameHeader] = h.String()
	nf[fieldNameLevel] = level
	nf[fieldNameMessage] = fmt.Sprintf(format, a...)
	traceId := tracing.TraceIdFromContext(ctx)
	if traceId != "" {
		nf[fieldNameTracing] = traceId
	}
	bs, _ := json.Marshal(nf)
	s := string(bs)
	select {
	case l.buffer <- s:
		return
	case <-time.After(time.Millisecond):
		fmt.Printf("[%v] drop log because log buffer is full %v\n", l.typ, s)
	}
}
