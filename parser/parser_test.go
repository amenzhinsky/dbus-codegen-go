package parser

import (
	"testing"
)

func TestParseSig(t *testing.T) {
	t.Parallel()
	for s, want := range map[string]string{
		"y":        "byte",
		"b":        "bool",
		"n":        "int16",
		"q":        "uint16",
		"i":        "int32",
		"u":        "uint32",
		"x":        "int64",
		"t":        "uint64",
		"d":        "float64",
		"h":        "dbus.UnixFD",
		"s":        "string",
		"o":        "dbus.ObjectPath",
		"v":        "dbus.Variant",
		"g":        "dbus.Signature",
		"ai":       "[]int32",
		"aai":      "[][]int32",
		"aaaa{sb}": "[][][]map[string]bool",
		"a{yv}":    "map[byte]dbus.Variant",
		"(ybv)":    "struct {V0 byte;V1 bool;V2 dbus.Variant}",
	} {
		if have := parseSig(s); have != want {
			t.Errorf("parseSig(%q) = %q, want %q", s, have, want)
		}
	}
}
