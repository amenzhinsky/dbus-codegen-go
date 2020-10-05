package main

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"os"
)

const DBUS_SERVICE_NAME = "org.example.Demo"
const DBUS_OBJECT_PATH = dbus.ObjectPath("/org/example/Demo")

type Demo struct {
}

func (d Demo) WelcomeMessage() (outputMessage string, err *dbus.Error) {
	return "Hello, this is example codebase for dbus-codegen-go", nil
}

func main() {
	demo := Demo{}

	dBusConnection, err := dbus.SessionBus()
	if err != nil {
		fmt.Errorf("error connecting to dbus session bus - %v", err)
		return
	}
	defer dBusConnection.Close()

	err = ExportOrg_Example_Demo(dBusConnection, DBUS_OBJECT_PATH, demo)
	if err != nil {
		fmt.Errorf("error in exporting %v", err)
	}

	reply, err := dBusConnection.RequestName(DBUS_SERVICE_NAME,
		dbus.NameFlagDoNotQueue)
	if err != nil {
		fmt.Errorf("error in registering request name %v", err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Errorf("%v name already taken", os.Stderr)
		os.Exit(1)
	}

	fmt.Printf("Listening on interface - %v and path %v ...\n", DBUS_SERVICE_NAME, DBUS_OBJECT_PATH)
	select {}
}
