package dbusgen

import (
	"reflect"
	"testing"
)

func TestSig(t *testing.T) {
	for s, want := range map[string]sig{
		"y":        {"byte"},
		"bb":       {"bool", "bool"},
		"nqiuxtd":  {"int16", "uint16", "int32", "uint32", "int64", "uint64", "float64"},
		"hsovg":    {"dbus.UnixFD", "string", "dbus.ObjectPath", "interface{}", "dbus.Signature"},
		"ai":       {"[]int32"},
		"aai":      {"[][]int32"},
		"aaaa{sb}": {"[][][]map[string]bool"},
		"a{yv}":    {"map[byte]interface{}"},
		"(yb)":     {"struct {byte;bool}"},
	} {
		if have := newSig(s); !reflect.DeepEqual(have, want) {
			t.Errorf("newSig(%q) = %v, want %v", s, have, want)
		}
	}
}
