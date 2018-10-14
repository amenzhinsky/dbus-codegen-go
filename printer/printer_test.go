package printer

import (
	"bytes"
	"testing"

	"github.com/amenzhinsky/godbus-codegen/token"
)

func TestPrint(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := Print(&buf, []*token.Interface{
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
