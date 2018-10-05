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

		want arg
	}{
		{"", "u", "arg", 8, arg{"arg8", "uint32"}},
		{"interface", "i", "var", 0, arg{"varInterface", "int32"}},
		{"my_varName", "o", "in", 1, arg{"myVarName", "dbus.ObjectPath"}},
		{"camel___case", "s", "in", 2, arg{"camelCase", "string"}},
		{"CamelCase", "s", "out", 3, arg{"camelCase", "string"}},
	} {
		if have := newArg(run.identifier, run.signature, run.prefix, run.i); have != run.want {
			t.Errorf("newArg(%q, %q, %q, %d) = %v, want %v",
				run.identifier, run.signature, run.prefix, run.i, have, run.want,
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
