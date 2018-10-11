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
}

func (b *buffer) writef(format string, v ...interface{}) (int, error) {
	return fmt.Fprintf(&b.buf, format, v...)
}

func (b *buffer) writeln(s ...string) {
	for i := 0; i < len(s); i++ {
		b.buf.WriteString(s[i])
	}
	b.buf.WriteByte('\n')
}

func (b *buffer) bytes() []byte {
	return b.buf.Bytes()
}
