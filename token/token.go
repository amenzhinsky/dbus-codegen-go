package token

import "github.com/godbus/dbus/v5/introspect"

// Interface is a D-Bus interface.
type Interface struct {
	Name         string
	Methods      []*Method
	Properties   []*Property
	Signals      []*Signal
	Annotations  []*Annotation
	RawInterface introspect.Interface
}

// Method is a D-Bus method.
type Method struct {
	Name        string
	In          []*Arg
	Out         []*Arg
	Annotations []*Annotation
}

// Property is a D-Bus property.
type Property struct {
	Name        string
	Arg         *Arg
	Read        bool
	Write       bool
	Annotations []*Annotation
}

// Signal is a D-Bus signal.
type Signal struct {
	Name        string
	Args        []*Arg
	Annotations []*Annotation
}

// Arg is an argument.
type Arg struct {
	Name string
	Type string
}

// Annotation is a D-Bus annotation.
type Annotation struct {
	Name  string
	Value string
}
