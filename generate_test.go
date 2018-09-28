package dbusgen

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGenerateFreedesktopDBus(t *testing.T) {
	g, err := New(WithPackageName("gengendbus"))
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadFile("testdata/test.xml")
	if err != nil {
		t.Fatal(err)
	}
	o, err := g.Parse(b)
	if err != nil {
		t.Fatal(err)
	}

	//tf, err := ioutil.TempFile("", "")
	tf, err := os.OpenFile("out/gen.go", os.O_WRONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		tf.Close()
		//os.Remove(tf.Name())
	}()

	if _, err = tf.Write(o); err != nil {
		t.Fatal(err)
	}
}
