package integration_test

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

// generated with: `dbus-codegen-go -xml -dest=X > testdata/X.xml`
var xmlFiles = []string{
	"testdata/net.connman.iwd.xml",
	"testdata/org.bluez.xml",
	"testdata/org.freedesktop.Accounts.xml",
	"testdata/org.freedesktop.ColorManager.xml",
	"testdata/org.freedesktop.DBus.xml",
	"testdata/org.freedesktop.hostname1.xml",
	"testdata/org.freedesktop.import1.xml",
	"testdata/org.freedesktop.locale1.xml",
	"testdata/org.freedesktop.login1.xml",
	"testdata/org.freedesktop.machine1.xml",
	"testdata/org.freedesktop.NetworkManager.xml",
	"testdata/org.freedesktop.resolve1.xml",
	"testdata/org.freedesktop.systemd1.xml",
	"testdata/org.freedesktop.timedate1.xml",
	"testdata/org.freedesktop.UDisks2.xml",
	"testdata/org.freedesktop.Upower.xml",
	"testdata/org.gnome.DisplayManager.xml",
}

func TestReproducibility(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	for _, f := range xmlFiles {
		f := f
		t.Run(f, func(t *testing.T) {
			t.Parallel()
			b1, err := generate(f)
			if err != nil {
				t.Fatal(err)
			}
			b2, err := generate(f)
			if err != nil {
				t.Fatal(err)
			}
			h1 := sha256.Sum256(b1)
			h2 := sha256.Sum256(b2)
			if h1 != h2 {
				t.Errorf("hashsums for %v are different over multiple runs", f)
			}
		})
	}
}

func TestItCompiles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	src := []byte(`package main
func main() {}`)

	for _, f := range xmlFiles {
		f := f
		t.Run(f, func(t *testing.T) {
			checkCompile(t, src, f)
			t.Run("camelize", func(t *testing.T) {
				checkCompile(t, src, "-camelize", f)
			})
			t.Run("server-only", func(t *testing.T) {
				checkCompile(t, src, "-server-only", f)
			})
			t.Run("client-only", func(t *testing.T) {
				checkCompile(t, src, "-client-only", f)
			})
		})
	}
}

func TestScenario(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	for _, tc := range [][]string{
		{"testdata/test_signal.go", "testdata/org.freedesktop.DBus.xml"},
		{"testdata/test_single_method.go", "testdata/org.freedesktop.DBus.xml"},
		{"testdata/test_properties.go", "testdata/org.freedesktop.DBus.xml"},
		{"testdata/test_server_export.go", "testdata/org.freedesktop.DBus.xml"},
		{"testdata/test_server_emit.go", "testdata/org.freedesktop.DBus.xml"},
	} {
		goFile, xmlFile := tc[0], tc[1]
		t.Run(goFile, func(t *testing.T) {
			b, err := ioutil.ReadFile(goFile)
			if err != nil {
				t.Fatal(err)
			}
			checkCompile(t, b, xmlFile)
		})
	}
}

func checkCompile(t *testing.T, src []byte, args ...string) {
	t.Helper()
	t.Parallel()
	b, err := generate(args...)
	if err != nil {
		t.Fatalf("generate(%v) error: %s", args, err)
	}
	if err := compile(b, src); err != nil {
		t.Fatalf("compile(%q, %v) error: %s", src, args, err)
	}
}

func generate(args ...string) ([]byte, error) {
	cmd := exec.Command("go",
		append([]string{"run", "../main.go", "-package=main"}, args...)...,
	)
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

func compile(gen, src []byte) error {
	temp, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)

	if err = ioutil.WriteFile(temp+"/gen.go", gen, 0644); err != nil {
		return err
	}
	if err = ioutil.WriteFile(temp+"/main.go", src, 0644); err != nil {
		return err
	}
	if out, err := exec.Command(
		"go", "run", temp+"/main.go", temp+"/gen.go",
	).CombinedOutput(); err != nil {
		return fmt.Errorf("compile error: %s", out)
	}
	return nil
}
