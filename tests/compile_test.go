package integration_test

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

var (
	binaryMu   sync.Mutex
	binaryFile string
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
	if binaryFile != "" { // no need to lock mutex
		os.Remove(binaryFile)
	}
}

// generated with: `dbus-codegen-go -xml -dest=X > testdata/X.xml`
var xmlFiles = []string{
	"testdata/net.connman.iwd.xml",
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

func xmlFilesAndAll() [][]string {
	all := make([][]string, 0, len(xmlFiles)+1)
	for _, file := range xmlFiles {
		all = append(all, []string{file})
	}
	return append(all, xmlFiles)
}

func TestHashSum(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	for _, tc := range xmlFilesAndAll() {
		files := tc
		t.Run(strings.Join(files, "/"), func(t *testing.T) {
			t.Parallel()
			h1 := sha256.Sum256(run(t, files...))
			h2 := sha256.Sum256(run(t, files...))
			if h1 != h2 {
				t.Errorf("hashsums for %v are different over multiple runs", files)
			}
		})
	}
}

func TestItCompiles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	for _, tc := range xmlFilesAndAll() {
		files := tc
		t.Run(strings.Join(files, "/"), func(t *testing.T) {
			t.Parallel()
			checkCompile(t, "testdata/test_it_compiles.gof", files)
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
	} {
		goFile, xmlFiles := tc[0], tc[1:]
		t.Run(goFile+"_"+strings.Join(xmlFiles, "-"), func(t *testing.T) {
			t.Parallel()
			checkCompile(t, goFile, xmlFiles)
		})
	}
}

func checkCompile(t *testing.T, goFile string, xmlFiles []string) {
	t.Helper()
	b := run(t, append([]string{"-package", "main"}, xmlFiles...)...)
	if err := compile(b, goFile); err != nil {
		t.Errorf("compile(%q, %v) error: %s", goFile, xmlFiles, err)
	}
}

// run runs the package binary with the given args.
// It compiles the package to a temporary file for possible further reuse,
// since `go run` takes much time for linking each time, TestMain cleans it up.
func run(t *testing.T, argv ...string) []byte {
	t.Helper()
	binaryMu.Lock()
	if binaryFile == "" {
		var err error
		binaryFile, err = build()
		if err != nil {
			t.Fatal(err)
		}
	}
	binaryMu.Unlock()

	b, err := exec.Command(binaryFile, argv...).CombinedOutput()
	if err != nil {
		t.Fatalf("generate error: %s, output: %s", err, string(b))
	}
	return b
}

func build() (string, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	f.Close()
	if err = os.Chmod(f.Name(), 0644); err != nil {
		return "", err
	}
	if b, err := exec.Command(
		"go", "build", "-ldflags=-s", "-o", f.Name(), "./..",
	).CombinedOutput(); err != nil {
		return "", fmt.Errorf("binary compilation error: %s", string(b))
	}
	return f.Name(), nil
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
