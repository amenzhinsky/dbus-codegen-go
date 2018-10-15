package integration_test

import (
	"crypto/sha256"
	"strings"
	"testing"
)

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
