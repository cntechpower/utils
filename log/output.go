package log

import (
	"encoding/json"
	"fmt"
	"time"
)

func logOutput(skip int, h *Header, level Level, format string, a ...interface{}) {
	file := ""
	line := 0
	if h.reportFileLine {
		file, line = getCaller(skip)
	}
	for _, l := range loggers {
		switch l.typ {
		case OutputTypeText:
			logOutputText(l, file, line, h, level, format, a...)
		case OutputTypeJson:
			logOutputStructured(l, file, line, h, level, format, a...)
		}
	}
}

func logOutputText(l *loggerWithConfig, file string, line int, h *Header, level Level, format string, a ...interface{}) {
	if h.reportFileLine {
		l.Println(fmt.Sprintf("[%s] <%s> |%s|%s| (%s:%v) %s",
			time.Now().Format("2006-01-02 15:04:05.000"), level, h, h.fields.String(),
			file, line, fmt.Sprintf(format, a...)))
	} else {
		l.Println(fmt.Sprintf("[%s] <%s> |%s|%s| %s",
			time.Now().Format("2006-01-02 15:04:05.000"), level, h, h.fields.String(),
			fmt.Sprintf(format, a...)))
	}

}

func logOutputStructured(l *loggerWithConfig, file string, line int, h *Header, level Level, format string, a ...interface{}) {
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
	bs, _ := json.Marshal(nf)
	l.Println(string(bs))
}
