package printer

import "testing"

func TestBuffer(t *testing.T) {
	t.Parallel()

	buf := newBuffer()
	buf.writeln("foo", "bar")
	buf.writef("%s", "baz")
	want := "foobar\nbaz"
	have, err := buf.bytes()
	if err != nil {
		t.Fatal(err)
	}
	if string(have) != want {
		t.Fatalf("buffer output = %q, want %q", string(have), want)
	}
}
