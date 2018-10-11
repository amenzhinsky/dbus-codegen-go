package integration_test

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/amenzhinsky/godbus-codegen/parser"
	"github.com/amenzhinsky/godbus-codegen/printer"
	"github.com/amenzhinsky/godbus-codegen/token"
)

func TestOutputHashSum(t *testing.T) {
	ifaces := parse(t, "org.freedesktop.DBus.xml")
	hash1 := sha256.New()
	hash2 := sha256.New()
	if err := printer.Print(hash1, ifaces); err != nil {
		t.Fatal(err)
	}
	if err := printer.Print(hash2, ifaces); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(hash1.Sum(nil), hash2.Sum(nil)) {
		t.Fatal("outputs have different hash sums")
	}
}

func TestCompile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	for _, run := range []struct {
		xml, gof string
	}{
		{"org.freedesktop.DBus.xml", "test_signal.gof"},
		{"org.freedesktop.DBus.xml", "test_single_method.gof"},
		{"org.freedesktop.DBus.xml", "test_properties.gof"},
	} {
		t.Run(run.gof, func(t *testing.T) {
			ifaces := parse(t, "org.freedesktop.DBus.xml")
			if err := compile(ifaces, run.gof); err != nil {
				t.Errorf("compile(%q, %q) error: %s", run.xml, run.gof, err)
			}
		})
	}
}

func parse(t *testing.T, xmlFile string) []*token.Interface {
	t.Helper()
	b, err := ioutil.ReadFile("testdata/" + xmlFile)
	if err != nil {
		t.Fatal(err)
	}
	ifaces, err := parser.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	return ifaces
}

func compile(ifaces []*token.Interface, goFile string) error {
	temp, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)

	f, err := os.OpenFile(temp+"/dbus.go", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if err = printer.Print(f, ifaces, printer.WithPackageName("main")); err != nil {
		return err
	}
	defer f.Close()

	if err := copyFile(temp+"/main.go", "testdata/"+goFile); err != nil {
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
