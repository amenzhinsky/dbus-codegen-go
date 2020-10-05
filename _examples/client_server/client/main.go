package main

import (
	"context"
	"fmt"
	"github.com/godbus/dbus/v5"
)

const DEST = "org.example.Demo"
const DBUS_OBJECT_PATH = dbus.ObjectPath("/org/example/Demo")

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	demoObject := conn.Object(DEST, DBUS_OBJECT_PATH)

	obj := NewOrg_Example_Demo(demoObject)
	message, err := obj.WelcomeMessage(context.Background())
	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Printf("Message over dBus : %v\n", message)
}
