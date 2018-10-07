package dbusgen

import (
	"strings"

	"github.com/godbus/dbus/introspect"
)

type signal struct {
	gtyp string
	name string
	args []arg
}

func parseSignals(gtyp string, sigs []introspect.Signal) []*signal {
	signals := make([]*signal, len(sigs))
	for i := 0; i < len(sigs); i++ {
		signals[i] = &signal{
			gtyp: gtyp + strings.Title(sigs[i].Name) + "Signal",
			name: sigs[i].Name,
			args: argsToGoArgs(sigs[i].Args, "", "prop", true),
		}
	}
	return signals
}
