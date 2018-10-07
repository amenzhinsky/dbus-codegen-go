package dbusgen

import (
	"strings"

	"github.com/godbus/dbus/introspect"
)

type signal struct {
	gtyp  string
	iface string
	name  string
	args  []arg
}

func parseSignals(gtyp, iface string, sigs []introspect.Signal) []*signal {
	signals := make([]*signal, len(sigs))
	for i := 0; i < len(sigs); i++ {
		signals[i] = &signal{
			gtyp:  gtyp + strings.Title(sigs[i].Name) + "Signal",
			iface: iface,
			name:  sigs[i].Name,
			args:  argsToGoArgs(sigs[i].Args, "", "prop", true),
		}
	}
	return signals
}
