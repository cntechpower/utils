package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	Init(WithStd(OutputTypeJson),
		WithFile(OutputTypeText, "/tmp/test-log.out"),
		WithEs("main.unit-test.log", "http://10.0.0.2:9200"),
		WithKafka("main.unit-test.log", "10.0.0.2:9093", "AsyncLogging"))
	SetDefaultFields(HostIpFields)
	defer Close()
	h := NewHeader("new").WithField("test", true)
	h.Infof("hello world, %v", "dujinyang")
}
