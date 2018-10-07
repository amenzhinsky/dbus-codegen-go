package printer

import (
	"bytes"
	"testing"

	"github.com/amenzhinsky/godbus-codegen/token"
)

func TestPrint(t *testing.T) {
	var buf bytes.Buffer
	if err := Print(&buf, "main", []*token.Interface{
		{
			Type:       "FooOrg",
			Name:       "foo.org",
			Methods:    []*token.Method{},
			Properties: []*token.Property{},
			Signals:    []*token.Signal{},
		},
	}); err != nil {
		t.Fatal(err)
	}

	// TODO: test s
}
