package trace

import (
	"fmt"
	"io"
)

// Tracerはコード内の出来事を記録できるオブジェクトを表すインターフェース
type Tracer interface {
	// ... で任意の型を何個でもという意味になる。
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}
