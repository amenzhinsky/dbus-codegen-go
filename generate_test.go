package dbusgen

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/amenzhinsky/godbus-codegen/parser"
)

func TestGenerate(t *testing.T) {
	for _, run := range []struct {
		xml, gof string
	}{
		{"org.freedesktop.DBus.xml", "get_id.gof"},
	} {
		if err := compile(run.xml, run.gof); err != nil {
			t.Errorf("compile(%q, %q) error: %s", run.xml, run.gof, err)
		}
	}
}

func compile(xmlFile, goFile string) error {
	g, err := New(WithPackageName("main"))
	if err != nil {
		return err
	}
	b, err := ioutil.ReadFile("testdata/" + xmlFile)
	if err != nil {
		return err
	}
	ifaces, err := parser.Parse(b)
	if err != nil {
		return err
	}
	o, err := g.Generate(ifaces...)
	if err != nil {
		return err
	}

	temp, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)

	if err := copyFile(temp+"/main.go", "testdata/"+goFile); err != nil {
		return err
	}
	if err := ioutil.WriteFile(temp+"/dbus.go", o, 0644); err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(temp); err != nil {
		return err
	}
	defer os.Chdir(cwd)

	out, err := exec.Command("go", "build", "-o", "a.out").CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation error: %s", out)
	}
	out, err = exec.Command("./a.out").CombinedOutput()
	if err != nil {
		return fmt.Errorf("executable error: %s", out)
	}
	return nil
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
