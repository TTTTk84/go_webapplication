package trace

import (
	"bytes"
	"io"
	"testing"
)
func TestNew(t *testing.T){
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り値がnilです")
	} else {
		tracer.Trace("こんにちは、traceパッケージ")
		if buf.String() != "こんにちは、traceパッケージ\n" {
			t.Errorf("'%s'という誤った文字列が出力されました", buf.String())
		}
	}
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
