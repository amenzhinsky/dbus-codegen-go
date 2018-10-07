package dbusgen

import (
	"testing"
)

func TestArg(t *testing.T) {
	for _, run := range []struct {
		identifier string
		signature  string
		prefix     string
		i          int
		export     bool

		want arg
	}{
		{"", "u", "arg", 8, false, arg{"arg8", "uint32"}},
		{"interface", "i", "var", 0, false, arg{"varInterface", "int32"}},
		{"my_varName", "o", "in", 1, false, arg{"myVarName", "dbus.ObjectPath"}},
		{"camel___case", "s", "in", 2, false, arg{"camelCase", "string"}},
		{"CamelCase", "s", "out", 3, false, arg{"camelCase", "string"}},
		{"exportVar", "s", "out", 4, true, arg{"ExportVar", "string"}},
	} {
		if have := newArg(run.identifier, run.signature, run.prefix, run.i, run.export); have != run.want {
			t.Errorf("newArg(%q, %q, %q, %d, %t) = %v, want %v",
				run.identifier, run.signature, run.prefix, run.i, run.export, have, run.want,
			)
		}
	}
}

func TestIfaceType(t *testing.T) {
	name, want := "org.freedesktop.DBus", "OrgFreedesktopDBus"
	if have := newIfaceType(name); have != want {
		t.Fatalf("newIfaceType(%q) = %q, want %q", name, have, want)
	}
}
