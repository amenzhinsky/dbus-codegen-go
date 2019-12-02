package main

import (
	"fmt"
	"os"

	"github.com/godbus/dbus/v5"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	defer conn.Close()

	if err = ExportOrg_Freedesktop_DBus_Introspectable(
		conn,
		"/org/test",
		&introspectable{},
	); err != nil {
		return err
	}
	return UnexportOrg_Freedesktop_DBus_Introspectable(conn, "/org/test")
}

type introspectable struct {
	UnimplementedOrg_Freedesktop_DBus_Introspectable
}
