package printer

import "testing"

func TestBuffer(t *testing.T) {
	buf := newBuffer()
	buf.writeln("foo", "bar")
	buf.writef("%s", "baz")
	want := "foobar\nbaz"
	if have := string(buf.bytes()); have != want {
		t.Fatalf("buffer output = %q, want %q", have, want)
	}
}
