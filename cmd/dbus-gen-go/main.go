package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/amenzhinsky/godbus-codegen"
	"github.com/godbus/dbus"
)

var (
	destFlag    string
	pathFlag    string
	ifaceFlag   string
	sessionFlag bool
	packageFlag string
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s [FLAG...]

Flags:
`, os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&destFlag, "dest", "", "destination name to introspect")
	flag.StringVar(&pathFlag, "path", "", "object path to introspect")
	flag.StringVar(&ifaceFlag, "iface", "", "generate only for the named interfaces, coma-separated")
	flag.BoolVar(&sessionFlag, "session", false, "connect to the session bus instead of the system")
	flag.StringVar(&packageFlag, "package", "dbusgen", "generated package name")
	flag.Parse()

	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(2)
	}
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	c, err := connect(sessionFlag)
	if err != nil {
		return err
	}
	defer c.Close()

	var b []byte
	if destFlag != "" || pathFlag != "" {
		if destFlag == "" {
			return errors.New("-dest is required for introspection")
		} else if pathFlag == "" {
			return errors.New("-path is required for introspection")
		}
		b, err = introspect(sessionFlag, destFlag, dbus.ObjectPath(pathFlag))
	} else {
		b, err = ioutil.ReadAll(os.Stdin)
	}
	if err != nil {
		return err
	}
	g, err := dbusgen.New(
		dbusgen.WithPackageName(packageFlag),
	)
	if err != nil {
		return err
	}
	output, err := g.Generate([][]byte{b}, split(ifaceFlag)...)
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func split(s string) []string {
	ss := make([]string, 0, strings.Count(s, ",")+1)
	for _, chunk := range strings.Split(s, ",") {
		if chunk = strings.Trim(chunk, " "); chunk != "" {
			ss = append(ss, chunk)
		}
	}
	return ss
}

func introspect(session bool, dest string, path dbus.ObjectPath) ([]byte, error) {
	conn, err := connect(session)
	if err != nil {
		return nil, err
	}
	var s string
	if err := conn.Object(dest, path).Call(
		"org.freedesktop.DBus.Introspectable.Introspect", 0,
	).Store(&s); err != nil {
		return nil, err
	}
	return []byte(s), nil
}

func connect(session bool) (*dbus.Conn, error) {
	if session {
		return dbus.SessionBus()
	}
	return dbus.SystemBus()
}
