package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/amenzhinsky/godbus-codegen/parser"
	"github.com/amenzhinsky/godbus-codegen/printer"
	"github.com/amenzhinsky/godbus-codegen/token"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

var (
	destFlag    string
	onlyFlag    []string
	exceptFlag  []string
	sessionFlag bool
	packageFlag string
)

type stringsFlag []string

func (ss *stringsFlag) String() string {
	return "[" + strings.Join(*ss, ", ") + "]"
}

func (ss *stringsFlag) Set(s string) error {
	if s = strings.Trim(s, " "); s == "" {
		return errors.New("string is empty")
	}
	*ss = append(*ss, s)
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s [FLAG...] [PATH...]

Generates code from DBus org.freedesktop.DBus.Introspectable format.
It introspects the named destination by providing -dest flag, reads 
files passed as the command arguments or reads it from STDIN.

Flags:
`, os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&destFlag, "dest", "", "destination name to introspect")
	flag.Var((*stringsFlag)(&onlyFlag), "only", "generate code only for the named interfaces")
	flag.Var((*stringsFlag)(&exceptFlag), "except", "skip the named interfaces")
	flag.BoolVar(&sessionFlag, "session", false, "connect to the session bus instead of the system")
	flag.StringVar(&packageFlag, "package", "dbusgen", "generated package name")
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	var ifaces []*token.Interface
	if destFlag != "" {
		if flag.NArg() > 0 {
			return errors.New("cannot combine arguments and -dest flag")
		}
		conn, err := connect(sessionFlag)
		if err != nil {
			return err
		}
		defer conn.Close()

		ifaces, err = parseDest(conn, destFlag, "/")
		if err != nil {
			return err
		}
	} else if flag.NArg() > 0 {
		for _, filename := range flag.Args() {
			b, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}
			chunk, err := parser.Parse(b)
			if err != nil {
				return err
			}
			ifaces = merge(ifaces, chunk)
		}
	} else {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		ifaces, err = parser.Parse(b)
		if err != nil {
			return err
		}
	}

	if len(onlyFlag) != 0 && len(exceptFlag) != 0 {
		return errors.New("cannot combine -only and -except flags")
	}
	filtered := make([]*token.Interface, 0, len(ifaces))
	for _, iface := range ifaces {
		if len(onlyFlag) == 0 && len(exceptFlag) == 0 ||
			len(onlyFlag) != 0 && includes(onlyFlag, iface.Name) ||
			len(exceptFlag) != 0 && !includes(exceptFlag, iface.Name) {
			filtered = append(filtered, iface)
		}
	}
	return printer.Print(os.Stdout, filtered, printer.WithPackageName(packageFlag))
}

func connect(session bool) (*dbus.Conn, error) {
	if session {
		return dbus.SessionBus()
	}
	return dbus.SystemBus()
}

func parseDest(conn *dbus.Conn, dest string, path dbus.ObjectPath) (
	[]*token.Interface, error,
) {
	ifaces, children, err := introspectPath(conn, destFlag, path)
	if err != nil {
		return nil, err
	}
	if path == "/" {
		path = ""
	}
	for _, child := range children {
		chunk, err := parseDest(conn, dest, path+dbus.ObjectPath("/"+child.Name))
		if err != nil {
			return nil, err
		}
		ifaces = merge(ifaces, chunk)
	}
	return ifaces, nil
}

func merge(curr, next []*token.Interface) []*token.Interface {
	for j := range next {
		var found bool
		for i := range curr {
			if curr[i].Name == next[j].Name {
				found = true
				break
			}
		}
		if !found {
			curr = append(curr, next[j])
		}
	}
	return curr
}

func includes(ss []string, s string) bool {
	for i := range ss {
		if ss[i] == s {
			return true
		}
	}
	return false
}

func introspectPath(conn *dbus.Conn, dest string, path dbus.ObjectPath) (
	[]*token.Interface, []introspect.Node, error,
) {
	var s string
	if err := conn.Object(dest, path).Call(
		"org.freedesktop.DBus.Introspectable.Introspect", 0,
	).Store(&s); err != nil {
		return nil, nil, err
	}
	var node introspect.Node
	if err := xml.Unmarshal([]byte(s), &node); err != nil {
		return nil, nil, err
	}
	ifaces, err := parser.ParseNode(&node)
	if err != nil {
		return nil, nil, err
	}
	return ifaces, node.Children, nil
}
