package parser

import (
	"encoding/xml"
	"go/token"
	"regexp"
	"strconv"
	"strings"

	"github.com/godbus/dbus/introspect"
)

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

func parseMethods(methods []introspect.Method) []*Method {
	list := make([]*Method, 0, len(methods))
	for _, method := range methods {
		list = append(list, &Method{
			Type: strings.Title(method.Name),
			Name: method.Name,
			In:   parseArgs(method.Args, "in", "arg", false),
			Out:  parseArgs(method.Args, "out", "ret", false),
		})
	}
	return list
}

type Property struct {
	Type    string
	Name    string
	Return  string
	Default string
	Read    bool
	Write   bool
}

func parseProperties(props []introspect.Property) []*Property {
	list := make([]*Property, 0, len(props))
	for _, prop := range props {
		list = append(list, &Property{
			Type:    strings.Title(prop.Name),
			Name:    prop.Name,
			Return:  parseSignature(prop.Type)[0],
			Default: signatureZeroValue(prop.Type),
			Read:    strings.Index(prop.Access, "read") >= 0,
			Write:   strings.Index(prop.Access, "write") >= 0,
		})
	}
	return list
}

type Signal struct {
	Type string
	Name string
	Args []*Arg
}

func parseSignals(typ string, sigs []introspect.Signal) []*Signal {
	list := make([]*Signal, 0, len(sigs))
	for _, sig := range sigs {
		list = append(list, &Signal{
			Type: typ + strings.Title(sig.Name) + "Signal",
			Name: sig.Name,
			Args: parseArgs(sig.Args, "", "prop", true),
		})
	}
	return list
}

type Arg struct {
	Name string
	Type string
}

func parseArgs(args []introspect.Arg, direction, prefix string, export bool) []*Arg {
	out := make([]*Arg, 0, len(args))
	for i := range args {
		if direction != "" && args[i].Direction != direction {
			continue
		}
		out = append(out, parseArg(args[i].Name, args[i].Type, prefix, len(out), export))
	}
	return out
}

var varRegexp = regexp.MustCompile("_+[a-zA-Z0-9]")

func parseArg(identifier, signature string, prefix string, i int, export bool) *Arg {
	var name string
	if identifier == "" {
		name = prefix + strconv.Itoa(i)
	} else if isKeyword(identifier) {
		name = prefix + strings.Title(identifier)
	} else {
		name = strings.ToLower(identifier[:1]) +
			varRegexp.ReplaceAllStringFunc(identifier[1:], func(s string) string {
				return strings.Title(strings.TrimLeft(s, "_"))
			})
	}
	if export {
		name = strings.Title(name)
	}
	return &Arg{Name: name, Type: parseSignature(signature)[0]}
}

func Parse(b []byte) ([]*Interface, error) {
	var node introspect.Node
	if err := xml.Unmarshal(b, &node); err != nil {
		return nil, err
	}
	var ifaces []*Interface
	for _, iface := range node.Interfaces {
		typ := ifaceToType(iface.Name)
		ifaces = append(ifaces, &Interface{
			Type:       typ,
			Name:       iface.Name,
			Methods:    parseMethods(iface.Methods),
			Properties: parseProperties(iface.Properties),
			Signals:    parseSignals(typ, iface.Signals),
		})
	}
	return ifaces, nil
}

var ifaceRegexp = regexp.MustCompile("\\.[a-zA-Z0-9]")

func ifaceToType(name string) string {
	name = strings.Title(name)
	if isKeyword(name) {
		return name
	}
	return ifaceRegexp.ReplaceAllStringFunc(name, func(s string) string {
		return strings.ToUpper(s[1:])
	})
}

func isKeyword(s string) bool {
	return token.Lookup(s).IsKeyword()
}
