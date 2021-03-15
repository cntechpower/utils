package log

import "testing"

func TestLog(t *testing.T) {
	Init(WithStd(OutputTypeJson), WithFile(OutputTypeText, "/tmp/test-log.out"))
	SetDefaultFields(HostIpFields)
	h := NewHeader("new").WithField("test", true)
	h.Infof("hello world, %v", "dujinyang")
}
