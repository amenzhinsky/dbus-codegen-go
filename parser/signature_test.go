package parser

import (
	"reflect"
	"testing"
)

func TestParseSignature(t *testing.T) {
	t.Parallel()
	for s, want := range map[string][]string{
		"y":        {"byte"},
		"bb":       {"bool", "bool"},
		"nqiuxtd":  {"int16", "uint16", "int32", "uint32", "int64", "uint64", "float64"},
		"hsovg":    {"dbus.UnixFD", "string", "dbus.ObjectPath", "dbus.Variant", "dbus.Signature"},
		"ai":       {"[]int32"},
		"aai":      {"[][]int32"},
		"aaaa{sb}": {"[][][]map[string]bool"},
		"a{yv}":    {"map[byte]dbus.Variant"},
		"(ybv)":    {"struct {v0 byte;v1 bool;v2 dbus.Variant}"},
	} {
		if have := parseSignature(s); !reflect.DeepEqual(have, want) {
			t.Errorf("parseSignature(%q) = %v, want %v", s, have, want)
		}
	}
}
