package token

// Interface is a D-Bus interface.
type Interface struct {
	Name       string
	Methods    []*Method
	Properties []*Property
	Signals    []*Signal
}

// Method is a D-Bus method.
type Method struct {
	Name string
	In   []*Arg
	Out  []*Arg
}

// Property is a D-Bus property.
type Property struct {
	Name  string
	Arg   *Arg
	Read  bool
	Write bool
}

// Signal is a D-Bus signal.
type Signal struct {
	Name string
	Args []*Arg
}

// Arg is an argument.
type Arg struct {
	Name string
	Type string
}
