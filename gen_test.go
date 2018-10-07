package dbusgen

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestGenerateFreedesktopDBus(t *testing.T) {
	g, err := New(WithPackageName("main"))
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadFile("testdata/test.xml")
	if err != nil {
		t.Fatal(err)
	}
	o, err := g.Generate(b)
	if err != nil {
		t.Fatal(err)
	}

	temp, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(temp)

	if err := copyFile(temp+"/main.go", "testdata/get_id.sample"); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(temp+"/dbus.go", o, 0644); err != nil {
		t.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(temp); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)

	out, err := exec.Command("go", "build", "-o", "a.out").CombinedOutput()
	if err != nil {
		t.Fatalf("compilation error: %s", out)
	}
	out, err = exec.Command("./a.out").CombinedOutput()
	if err != nil {
		t.Fatalf("executable error: %s", out)
	}
}

func copyFile(dst, src string) error {
	d, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer d.Close()

	s, err := os.OpenFile(src, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = io.Copy(d, s)
	return err
}
