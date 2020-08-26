# `org.freedesktop.DBus.ObjectManager` client example

To generate client-side code using `dbus-codegen-go` as a `go:generate` directive run: 

```bash
go generate
```

Now you can list and watch an object manager initial state and changes, for example `org.bluez`: 

```bash
go run . -system org.bluez
```

In a parallel console you can trigger events activity:

```bash
bluetoothctl
power on
scan on
```
