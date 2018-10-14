# dbus-codegen-go

Generates golang code for the given D-Bus introspection input.

## Overview

TODO

## Usage

Apart of library code, the repo provides `dbus-gen-go` CLI tool, that can be installed with the following command and make sure that `$GOPATH/bin` is in your `$PATH`:

```
go get -u github.com/amenzhinsky/godbus-codegen/cmd/dbus-gen-go
```

It can introspect destinations on a running DBus: 

```
dbus-gen-go -dest org.freedesktop.systemd1
```

It reads the provided xml files or if none given it reads from the stdin:

```
dbus-send --system \
	--type=method_call \
	--print-reply=literal \
	--dest=org.freedesktop.systemd1 \
	/org/freedesktop/systemd1 \
	org.freedesktop.DBus.Introspectable.Introspect \
	> org.freedesktop.systemd1.xml

dbus-gen-go org.freedesktop.systemd1.xml
cat org.freedesktop.systemd1.xml | dbus-gen-go
```

Here's an example of a bit more advanced usage, where we're changing the generated code's package name and narrow down the introspected interface to just two that we need:

```
dbus-gen-go \
	-dest org.freedesktop.systemd1 \
	-package systemd \
	-only org.freedesktop.systemd1.Manager \
	-only org.freedesktop.systemd1.Service
```

## Troubleshooting

### parse error: ...

The generated output by `printer` package cannot be parsed by gofmt, disable it by `-gofmt=false` if you're using CLI or `printer.WithGofmt(false)` ir you're using the repo as a library and inspect the result or fill in an issue with input xml files and the generated code.
