package parser

import "testing"

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

		want Arg
	}{
		{"", "u", "arg", 8, false, Arg{"arg8", "uint32"}},
		{"interface", "i", "var", 0, false, Arg{"varInterface", "int32"}},
		{"my_varName", "o", "in", 1, false, Arg{"myVarName", "dbus.ObjectPath"}},
		{"camel___case", "s", "in", 2, false, Arg{"camelCase", "string"}},
		{"CamelCase", "s", "out", 3, false, Arg{"camelCase", "string"}},
		{"exportVar", "s", "out", 4, true, Arg{"ExportVar", "string"}},
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
