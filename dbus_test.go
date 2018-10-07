package dbusgen

import (
	"reflect"
	"testing"

	"github.com/godbus/dbus/introspect"
)

func TestParseSignals(t *testing.T) {
	for _, run := range []struct {
		gtyp  string
		iface string
		sigs  []introspect.Signal
		want  []*signal
	}{
		{
			"OrgBluez", "org.bluez", []introspect.Signal{
				{Name: "ValueChanged", Args: []introspect.Arg{
					{
						Name: "prop",
						Type: "u",
					},
				}},
			},
			[]*signal{
				{"OrgBluezValueChangedSignal", "org.bluez", "ValueChanged", []arg{
					{"Prop", "uint32"},
				}},
			},
		},
	} {
		if have := parseSignals(run.gtyp, run.iface, run.sigs); !reflect.DeepEqual(have, run.want) {
			t.Errorf("parseSignals(%q, %v) = %v, want %v", run.gtyp, run.sigs, have, run.want)
		}
	}
}
