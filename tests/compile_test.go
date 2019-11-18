package integration_test

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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

func TestHashSum(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	for _, file := range xmlFiles {
		t.Run(file, func(t *testing.T) {
			b1, err := generate(file)
			if err != nil {
				t.Fatal(err)
			}
			b2, err := generate(file)
			if err != nil {
				t.Fatal(err)
			}
			h1 := sha256.Sum256(b1)
			h2 := sha256.Sum256(b2)
			if h1 != h2 {
				t.Errorf("hashsums for %v are different over multiple runs", file)
			}
		})
	}
}

func TestItCompiles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	for _, file := range xmlFiles {
		t.Run(file, func(t *testing.T) {
			checkCompile(t, "testdata/test_it_compiles.gof", file)
		})
	}
}

func TestCompile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	for _, tc := range [][]string{
		{"testdata/test_signal.gof", "testdata/org.freedesktop.DBus.xml"},
		{"testdata/test_single_method.gof", "testdata/org.freedesktop.DBus.xml"},
		{"testdata/test_properties.gof", "testdata/org.freedesktop.DBus.xml"},
		{"testdata/test_server_export.gof", "testdata/org.freedesktop.DBus.xml"},
		{"testdata/test_server_emit.gof", "testdata/org.freedesktop.DBus.xml"},
	} {
		goFile, xmlFile := tc[0], tc[1]
		t.Run(goFile, func(t *testing.T) {
			checkCompile(t, goFile, xmlFile)
		})
	}
}

func checkCompile(t *testing.T, goFile string, xmlFile string) {
	t.Helper()
	b, err := generate(xmlFile)
	if err != nil {
		t.Fatalf("generate(%q) error: %s", xmlFile, err)
	}
	if err := compile(b, goFile); err != nil {
		t.Fatalf("compile(%q, %v) error: %s", goFile, xmlFile, err)
	}
}

func generate(args ...string) ([]byte, error) {
	cmd := exec.Command("go",
		append([]string{"run", "../main.go", "-package=main"}, args...)...,
	)
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

func compile(b []byte, goFile string) error {
	temp, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)

	if err = ioutil.WriteFile(temp+"/gen.go", b, 0644); err != nil {
		return err
	}
	path, err := filepath.Abs(goFile)
	if err != nil {
		return err
	}
	if err = os.Symlink(path, temp+"/main.go"); err != nil {
		return err
	}
	if out, err := exec.Command(
		"go", "run", temp+"/main.go", temp+"/gen.go",
	).CombinedOutput(); err != nil {
		return fmt.Errorf("compile error: %s", out)
	}
	return nil
}
