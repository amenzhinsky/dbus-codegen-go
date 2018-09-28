package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/amenzhinsky/godbus-codegen"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

var (
	destFlag   string
	pathFlag   string
	ifaceFlag  string
	systemFlag bool
)

func main() {
	flag.StringVar(&destFlag, "dest", "", "destination name")
	flag.StringVar(&pathFlag, "path", "", "object path")
	flag.StringVar(&ifaceFlag, "iface", "", "interface to inspect")
	flag.BoolVar(&systemFlag, "system", false, "connect to the system bus instead of the session")
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	c, err := connect(systemFlag)
	if err != nil {
		return err
	}
	defer c.Close()

	if destFlag == "" {
		return errors.New("-dest is mandatory")
	} else if pathFlag == "" {
		return errors.New("-path is mandatory")
	} else if ifaceFlag == "" {
		return errors.New("-iface is mandatory")
	}

	n, err := introspect.Call(c.Object(destFlag, dbus.ObjectPath(pathFlag)))
	if err != nil {
		return err
	}

	g, err := dbusgen.New()
	if err != nil {
		return err
	}
	for _, iface := range n.Interfaces {
		if iface.Name == ifaceFlag {
			b, err := g.Generate(&iface)
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil
		}
	}
	return fmt.Errorf("interface %q not found", ifaceFlag)
}

func connect(system bool) (*dbus.Conn, error) {
	if system {
		return dbus.SystemBus()
	}
	return dbus.SessionBus()
}
