package token

type Interface struct {
	Type       string
	Name       string
	Methods    []*Method
	Properties []*Property
	Signals    []*Signal
}

type Method struct {
	Type string
	Name string
	In   []*Arg
	Out  []*Arg
}

type Property struct {
	Type    string
	Name    string
	Return  string
	Default string
	Read    bool
	Write   bool
}

type Signal struct {
	Type string
	Name string
	Args []*Arg
}

type Arg struct {
	Name string
	Type string
}
