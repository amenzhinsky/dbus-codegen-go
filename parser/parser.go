package parser

import (
	"encoding/xml"
	gotoken "go/token"
	"regexp"
	"strconv"
	"strings"

	"github.com/amenzhinsky/godbus-codegen/token"
	"github.com/godbus/dbus/introspect"
)

func parseMethods(methods []introspect.Method) []*token.Method {
	list := make([]*token.Method, 0, len(methods))
	for _, method := range methods {
		list = append(list, &token.Method{
			Type: strings.Title(method.Name),
			Name: method.Name,
			In:   parseArgs(method.Args, "in", "arg", false),
			Out:  parseArgs(method.Args, "out", "ret", false),
		})
	}
	return list
}

func parseProperties(props []introspect.Property) []*token.Property {
	list := make([]*token.Property, 0, len(props))
	for _, prop := range props {
		list = append(list, &token.Property{
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

func parseSignals(typ string, sigs []introspect.Signal) []*token.Signal {
	list := make([]*token.Signal, 0, len(sigs))
	for _, sig := range sigs {
		list = append(list, &token.Signal{
			Type: typ + strings.Title(sig.Name) + "Signal",
			Name: sig.Name,
			Args: parseArgs(sig.Args, "", "prop", true),
		})
	}
	return list
}

func parseArgs(args []introspect.Arg, direction, prefix string, export bool) []*token.Arg {
	out := make([]*token.Arg, 0, len(args))
	for i := range args {
		if direction != "" && args[i].Direction != direction {
			continue
		}
		out = append(out, parseArg(args[i].Name, args[i].Type, prefix, len(out), export))
	}
	return out
}

var varRegexp = regexp.MustCompile("_+[a-zA-Z0-9]")

func parseArg(identifier, signature string, prefix string, i int, export bool) *token.Arg {
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
	return &token.Arg{Name: name, Type: parseSignature(signature)[0]}
}

func Parse(b []byte) ([]*token.Interface, error) {
	var node introspect.Node
	if err := xml.Unmarshal(b, &node); err != nil {
		return nil, err
	}
	var ifaces []*token.Interface
	for _, iface := range node.Interfaces {
		typ := ifaceToType(iface.Name)
		ifaces = append(ifaces, &token.Interface{
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
	return gotoken.Lookup(s).IsKeyword()
}
