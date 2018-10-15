# dbus-codegen-go

Takes the [D-Bus Introspection Data Format](https://dbus.freedesktop.org/doc/dbus-specification.html#introspection-format) and generates Go code.

The project depends only on [github.com/godbus/dbus](https://github.com/godbus/dbus) module and cannot be used separately because it operates its data types.

API may change until 1.0.0, so please vendor the source code if you want to use this tool.

## Overview

XML like this:

```
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
1. `IToA` method attached to the structure: `(*My_Awesome_Interface) IToA(int64) (string, error)`.
1. `Powered` property getter and setter: `(*My_Awesome_Interface) GetPowered() (bool, error)` and `(*My_Awesome_Interface) SetPowered(bool) error`.
1. `My_Awesome_Interface_SomethingHappenedSignal` for typed access to signal body attributes and `LookupSignal(*dbus.Signal)` function for easy conversion:
  ```
  sig := LookupSignal(signal).(*My_Awesome_Interface_SomethingHappenedSignal)
  	fmt.Printf("%s %s", sig.Body().ObjectPath, sig.Body().What)
  ```  

## Installation

You can install it with `go get` withing `$GOPATH` or a module (**go1.11** has issues with binaries installation outside of a module):

```
go get -u github.com/amenzhinsky/dbus-codegen-go/cmd/dbus-codegen-go
```

Or clone the repo and build it manually:

```
git clone https://github.com/amenzhinsky/dbus-codegen-go.git
cd dbus-codegen-go
go build
```

## Usage

The program reads files provided as command-line arguments or reads the stdin if none given:

```
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

Apart of reading existing files it can introspect destinations on locally running D-Bus: 

```
dbus-codegen-go -dest=org.freedesktop.systemd1
```

You may also want to safe the introspection file that combines all interfaces in the tree on some system for further reuse. For that simply add `-xml` flag:

```
dbus-codegen-go -xml -dest=org.freedesktop.systemd1
```

Here's an example of a bit more advanced usage, where we're changing the generated code's package name and narrow down the introspected interface to just two that we need:

```
dbus-codegen-go \
	-dest=org.freedesktop.systemd1 \
	-package=systemd \
	-only=org.freedesktop.systemd1.Manager \
	-only=org.freedesktop.systemd1.Service
```

## Testing

To test the package simply run:

```
go test ./...
```

`tests` directory contains integration tests that technically compile the package's binary and use it for code generation and further compilation with test file scenarios, so make sure that `go` is in your `$PATH`.

## Troubleshooting

### parse error: ...

The generated output by `printer` package cannot be parsed by gofmt and that is the package issue, disable it with `-gofmt=false` and inspect the result or create an issue with input xml files and the generated code.

## TODO

- server side code generation
- add coding examples
- sophisticated tests
- more printer options, like trimming destination names or removing underscores from struct names

## Contributing

All contributions are welcome, just create an issue or issue a pull request.
