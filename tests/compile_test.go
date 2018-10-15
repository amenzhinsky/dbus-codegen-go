package integration_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

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
