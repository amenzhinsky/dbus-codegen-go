package integration_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

func xmlFilesAndAll() [][]string {
	all := make([][]string, 0, len(xmlFiles)+1)
	for _, file := range xmlFiles {
		all = append(all, []string{file})
	}
	return append(all, xmlFiles)
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
