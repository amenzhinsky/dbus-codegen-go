package parser

import (
	"testing"

	"github.com/amenzhinsky/godbus-codegen/token"
)

func TestIfaceType(t *testing.T) {
	name, want := "org.freedesktop.DBus", "OrgFreedesktopDBus"
	if have := ifaceToType(name); have != want {
		t.Fatalf("newIfaceType(%q) = %q, want %q", name, have, want)
	}
}

func TestParseArg(t *testing.T) {
	for _, run := range []struct {
		identifier string
		signature  string
		prefix     string
		i          int
		export     bool

		want token.Arg
	}{
		{"", "u", "arg", 8, false, token.Arg{Name: "arg8", Type: "uint32"}},
		{"interface", "i", "var", 0, false, token.Arg{Name: "varInterface", Type: "int32"}},
		{"my_varName", "o", "in", 1, false, token.Arg{Name: "myVarName", Type: "dbus.ObjectPath"}},
		{"camel___case", "s", "in", 2, false, token.Arg{Name: "camelCase", Type: "string"}},
		{"CamelCase", "s", "out", 3, false, token.Arg{Name: "camelCase", Type: "string"}},
		{"exportVar", "s", "out", 4, true, token.Arg{Name: "ExportVar", Type: "string"}},
	} {
		if have := parseArg(
			run.identifier, run.signature, run.prefix, run.i, run.export,
		); *have != run.want {
			t.Errorf("parseArg(%q, %q, %q, %d, %t) = %v, want %v",
				run.identifier, run.signature, run.prefix, run.i, run.export, *have, run.want,
			)
		}
	}
}
