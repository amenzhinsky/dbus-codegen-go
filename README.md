# dbus-codegen-go

Takes the [D-Bus Introspection Data Format](https://dbus.freedesktop.org/doc/dbus-specification.html#introspection-format) and generates Go code.

The project depends only on [github.com/godbus/dbus](https://github.com/godbus/dbus) module and cannot be used separately because it operates its data types.

API may change until 1.0.0, so please vendor the source code if you want to use this tool.

## Overview

XML like this:

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

The tool will generate the following data structures:

1. Structure `My_Awesome_Interface`, that can be created with `NewMy_Awesome_Interface(object)` or via `InterfaceLookup(object, "my.awesome.interface").(*My_Awesome_Interface)` in case you have more than one interface.
1. `IToA` method attached to the structure: `(*My_Awesome_Interface) IToA(context.Context, int64) (string, error)`.
1. `Powered` property getter and setter: `(*My_Awesome_Interface) GetPowered(context.Context) (bool, error)` and `(*My_Awesome_Interface) SetPowered(context.Context, bool) error`.
1. `My_Awesome_Interface_SomethingHappenedSignal` for typed access to signal body attributes, `LookupSignal(*dbus.Signal) (Signal, error)` and `AddMatchRule(*dbus.Signal) string` helper functions, see usage in the [examples](#examples) section.
    
1. Annotations added to interfaces, methods, properties and signals as comments.

## Installation

You can install it with `go get` withing `$GOPATH` or a module (**go1.11** has issues with binaries installation outside of a module):

```bash
go get -u github.com/amenzhinsky/dbus-codegen-go/cmd/dbus-codegen-go
```

Or clone the repo and build it manually:

```bash
git clone https://github.com/amenzhinsky/dbus-codegen-go.git
cd dbus-codegen-go
go build
```

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
	-only=org.freedesktop.systemd1.Manager \
	-only=org.freedesktop.systemd1.Service \
	-prefix=org.freedesktop.systemd1
```

## Examples

The following example subscribes to all `PropertyChanged` signals from `org.freedesktop.systemd1` destination.

The generated code is generated with:

```bash
dbus-codegen-go \
	-package=main \
	-dest=org.freedesktop.DBus \
	-dest=org.freedesktop.systemd1 
```

```go
package main

import (
	"fmt"
	"os"

	"github.com/godbus/dbus"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	defer conn.Close()

	sigc := make(chan *dbus.Signal, 1)
	conn.Signal(sigc)
	defer conn.RemoveSignal(sigc)

	bus := NewOrg_Freedesktop_DBus(conn.BusObject())
	if err := bus.AddMatch(
		AddMatchRule((*Org_Freedesktop_DBus_Properties_PropertiesChangedSignal)(nil)) +
			",sender='org.freedesktop.systemd1'",
	); err != nil {
		return err
	}
	for s := range sigc {
		sig, err := LookupSignal(s)
		if err != nil {
			return err
		}
		switch v := sig.(type) {
		case *Org_Freedesktop_DBus_Properties_PropertiesChangedSignal:
			fmt.Printf("%s %s: %v\n", v.Path(), v.Body.Interface, v.Body.ChangedProperties)
		}
	}
	return nil
}
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
- server side code generation
- add coding examples
- sophisticated tests
- more printer options, like using Ugly_Case and CamelCase in interface names

## Contributing

All contributions are welcome, just create an issue or issue a pull request.
