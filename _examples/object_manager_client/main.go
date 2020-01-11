//go:generate dbus-codegen-go -package=main -prefix=org.freedesktop.DBus -camelize -output=om.go -client-only om.xml
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/godbus/dbus/v5"
)

var (
	pathFlag   dbus.ObjectPath
	systemFlag bool
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [option...] DESTINATION\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar((*string)(&pathFlag), "path", "/", "bus object `path`")
	flag.BoolVar(&systemFlag, "system", false, "connect to the session bus")
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	if err := run(flag.Arg(0)); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run(dest string) error {
	conn, err := connect(systemFlag)
	if err != nil {
		return err
	}
	defer conn.Close()

	om := NewObjectManager(conn.Object(dest, pathFlag))
	objects, err := om.GetManagedObjects(context.Background())
	if err != nil {
		return err
	}
	for path, ifaces := range objects {
		fmt.Println(path)
		for iface := range ifaces {
			fmt.Printf("  %s\n", iface)
		}
		fmt.Println()
	}

	sigc := make(chan *dbus.Signal, 1)
	conn.Signal(sigc)
	defer conn.RemoveSignal(sigc)
	for _, sig := range []Signal{
		(*ObjectManagerInterfacesAddedSignal)(nil),
		(*ObjectManagerInterfacesRemovedSignal)(nil),
	} {
		if err = AddMatchSignal(conn, sig, dbus.WithMatchSender(dest)); err != nil {
			return err
		}
	}

	ossigc := make(chan os.Signal, 1)
	signal.Notify(ossigc, os.Interrupt)

	for {
		select {
		case s := <-sigc:
			v, err := LookupSignal(s)
			if err != nil {
				if err == ErrUnknownSignal {
					continue
				}
			}

			switch sig := v.(type) {
			case *ObjectManagerInterfacesAddedSignal:
				fmt.Printf("%s NEW\n", sig.Body.Object)
				for iface := range sig.Body.Interfaces {
					fmt.Printf("  %s\n", iface)
				}
				fmt.Println()
			case *ObjectManagerInterfacesRemovedSignal:
				fmt.Printf("%s DEL\n", sig.Body.Object)
				for _, iface := range sig.Body.Interfaces {
					fmt.Printf("  %s\n", iface)
				}
				fmt.Println()
			default:
				panic("should never get here")
			}
		case <-ossigc:
			signal.Reset()
			return nil
		}
	}
}

func connect(system bool) (*dbus.Conn, error) {
	if system {
		return dbus.SystemBus()
	}
	return dbus.SessionBus()
}
