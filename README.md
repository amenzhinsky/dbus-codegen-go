# dbus-codegen-go

[D-Bus Introspection Data Format](https://dbus.freedesktop.org/doc/dbus-specification.html#introspection-format) Go code generator.

The project depends only on [github.com/godbus/dbus](https://github.com/godbus/dbus) module and cannot be used separately because it operates its data types.

CLI and generated code API is a subject to change until `v1.0.0`.

## Overview

The tool generates two types of code-bases: server and client (or can be limited to just one of them by using `-client-only` or `-server-only` options), so if we take use the following XML:

```xml
<node>
	<interface name="my.awesome.interface">
		<method name="IToA">
			<arg name="In" type="x" direction="in" />
			<arg name="Out" type="s" direction="out" />
		</method>
		<property name="Powered" type="b" access="readwrite" />
		<signal name="SomethingHappened">
			<arg name="object_path" type="o" />
			<arg name="what" type="s" />
		</signal>
	</interface>
</node>
```

### Client

The program generates statically compiled bindings that are used to create objects that wrap around `dbus.BusObject`:

```go
obj := NewMy_Awesome_Interface(conn.Object("my.awesome.service", "/my/awesome/service"))
```

Now we can call its methods:

```go
s, err := obj.IToA(context.Background(), 666)
if err != nil {
    return err
}
fmt.Printf("itoa(%d) = %s", 666, s)
```

Retrieve or set property values:

```go
powered, err := obj.GetPowered(context.Background())
if err != nil {
    return err
}
fmt.Printf("powered = %t", powered)

if err := o.SetPowered(context.Background(), true); err != nil {
    return err
}
```

Handling signals requires type assertions:

```go
sigt := (*My_Awesome_Interface_SomethingHappenedSignal)(nil)
if err := AddMatchSignal(conn, sigt); err != nil {
    return err
}
defer RemoveMatchSignal(conn, sigt)

sigc := make(chan *dbus.Signal, 1)
conn.Signal(sigc)
for sig := range sigc {
    s, err := LookupSignal(sig)
    if err != nil {
        if err == ErrUnknownSignal {
            continue
        }
        return err
    }

    switch typed := s.(type) {
    case *My_Awesome_Interface_SomethingHappenedSignal:
        fmt.Printf("%s happened at %s", typed.Body.What, typed.Body.ObjectPath)
    }
}
```

### Server

Now we can implement the server interface, all interfaces are postfixed with `er`, in this example it's `My_Awesome_Interfaceer`.

It's recommended to embed the corresponding generated unimplemented structure for forward compatible implementations:

```go
type server struct {
	*UnimplementedMy_Awesome_Interface
}

func (s *server) IToA(in int64) (string, *dbus.Error) {
	return strconv.Itoa(int(in)), nil
}
```

And now we can export the implementation:

```go
if err := ExportMy_Awesome_Interface(conn, "/my/awesome/service", &srv{}); err != nil {
	return err
}
```

Emitting signals is done reusing the same structures generated for the client side:

```go
if err := Emit(conn, &My_Awesome_Interface_SomethingHappenedSignal{
    Path: "/org/my/iface",
    Body: &My_Awesome_Interface_SomethingHappenedSignalBody{
        ObjectPath: "/org/obj",
        What:       "something terrible",
    },
}); err != nil {
	return err
}
```

## Installation

You can install it with `go get`:

```bash
GO111MODULE=on go get -u github.com/amenzhinsky/dbus-codegen-go
```

Or clone the repo and build it manually:

```bash
git clone https://github.com/amenzhinsky/dbus-codegen-go.git .
go install
```

Make sure `$(go env GOPATH)/bin` is in your `$PATH`.

## Usage

The program treats command-line arguments as paths to XML files or reads out **stdin** if none given:

```bash
dbus-send --system \
	--type=method_call \
	--print-reply=literal \
	--dest=org.freedesktop.systemd1 \
	/org/freedesktop/systemd1 \
	org.freedesktop.DBus.Introspectable.Introspect \
	> org.freedesktop.systemd1.xml

dbus-codegen-go org.freedesktop.systemd1.xml
dbus-codegen-go < org.freedesktop.systemd1.xml
```

Apart of reading existing files it can introspect real D-Bus destinations recursively: 

```bash
dbus-codegen-go -dest=org.freedesktop.systemd1
```

You may also want to safe the introspection file that combines all interfaces in the tree on some system for further reuse. For that simply add `-xml` flag:

```bash
dbus-codegen-go -xml -dest=org.freedesktop.systemd1
```

Here's an example of a bit more advanced usage, where we're changing the generated code's package name and narrow down the introspected interfaces to just two we need, plus we're trimming `org.freedesktop` prefix to shorten generated structure names:

```bash
dbus-codegen-go \
	-dest=org.freedesktop.systemd1 \
	-package=systemd \
	-camelize \
	-only=org.freedesktop.systemd1.Manager \
	-only=org.freedesktop.systemd1.Service \
	-prefix=org.freedesktop.systemd1
```

## Testing

To test the package simply run:

```bash
go test ./...
```

`tests` directory contains integration tests that technically compile the package's binary and use it for code generation and further compilation with test file scenarios, so make sure that `go` is in your `$PATH`.

## Troubleshooting

### parse error: ...

The generated output by `printer` package cannot be parsed by gofmt and that is the package issue, disable it with `-gofmt=false` and inspect the result or create an issue with input xml files and the generated code.

## TODO

- name conflicts resolver
- add coding examples
- sophisticated tests
- handle no-reply calls

## Contributing

All contributions are welcome, just create an issue or issue a pull request.
