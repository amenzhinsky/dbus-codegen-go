package printer

import (
	"bytes"
	"fmt"
)

func newBuffer() *buffer {
	return &buffer{}
}

type buffer struct {
	buf bytes.Buffer
	err error
}

func (b *buffer) writef(format string, v ...interface{}) {
	if _, err := fmt.Fprintf(&b.buf, format, v...); err != nil && b.err != nil {
		b.err = err
	}
}

func (b *buffer) writeln(s ...string) {
	for i := 0; i < len(s); i++ {
		b.buf.WriteString(s[i])
	}
	b.buf.WriteByte('\n')
}

func (b *buffer) bytes() ([]byte, error) {
	return b.buf.Bytes(), b.err
}
