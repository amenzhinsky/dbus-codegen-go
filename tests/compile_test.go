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
	for _, tc := range [][2]string{
		{"test_signal.gof", "org.freedesktop.DBus.xml"},
		{"test_single_method.gof", "org.freedesktop.DBus.xml"},
		{"test_properties.gof", "org.freedesktop.DBus.xml"},

		{"test_it_compiles.gof", "net.connman.iwd.xml"},
		{"test_it_compiles.gof", "org.freedesktop.DBus.xml"},
		{"test_it_compiles.gof", "org.freedesktop.Accounts.xml"},
		{"test_it_compiles.gof", "org.freedesktop.UDisks2.xml"},
		{"test_it_compiles.gof", "org.freedesktop.systemd1.xml"},
		{"test_it_compiles.gof", "org.freedesktop.NetworkManager.xml"},
	} {
		tc := tc
		t.Run(tc[0], func(t *testing.T) {
			t.Parallel()
			ifaces := parse(t, tc[1])
			if err := compile(ifaces, tc[0]); err != nil {
				t.Errorf("compile(%q, %q) error: %s", tc[0], tc[1], err)
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

	f, err := os.OpenFile(temp+"/gen.go", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
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
	out, err := exec.Command("go", "run", temp+"/main.go", temp+"/gen.go").CombinedOutput()
	if err != nil {
		return fmt.Errorf("compile error: %s", out)
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
