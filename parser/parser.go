package parser

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/amenzhinsky/dbus-codegen-go/token"
	"github.com/godbus/dbus/introspect"
)

func Parse(b []byte) ([]*token.Interface, error) {
	var node introspect.Node
	if err := xml.Unmarshal(b, &node); err != nil {
		return nil, err
	}
	return ParseNode(&node)
}

func ParseNode(node *introspect.Node) ([]*token.Interface, error) {
	if node == nil {
		panic("node is nil")
	}
	var ifaces []*token.Interface
	for _, iface := range node.Interfaces {
		ifaces = append(ifaces, &token.Interface{
			Name:       iface.Name,
			Methods:    parseMethods(iface.Methods),
			Properties: parseProperties(iface.Properties),
			Signals:    parseSignals(iface.Signals),
		})
	}
	return ifaces, nil
}

func parseMethods(methods []introspect.Method) []*token.Method {
	list := make([]*token.Method, 0, len(methods))
	for _, method := range methods {
		list = append(list, &token.Method{
			Name: method.Name,
			In:   parseArgs(method.Args, "in"),
			Out:  parseArgs(method.Args, "out"),
		})
	}
	return list
}

func parseProperties(props []introspect.Property) []*token.Property {
	properties := make([]*token.Property, 0, len(props))
	for _, prop := range props {
		properties = append(properties, &token.Property{
			Name:  prop.Name,
			Arg:   parseArg(prop.Name, prop.Type),
			Read:  strings.Index(prop.Access, "read") != -1,
			Write: strings.Index(prop.Access, "write") != -1,
		})
	}
	return properties
}

func parseSignals(sigs []introspect.Signal) []*token.Signal {
	signals := make([]*token.Signal, 0, len(sigs))
	for _, sig := range sigs {
		signals = append(signals, &token.Signal{
			Name: sig.Name,
			Args: parseArgs(sig.Args, ""),
		})
	}
	return signals
}

func parseArgs(args []introspect.Arg, direction string) []*token.Arg {
	out := make([]*token.Arg, 0, len(args))
	for i := range args {
		if direction != "" && args[i].Direction != direction {
			continue
		}
		out = append(out, parseArg(args[i].Name, args[i].Type))
	}
	return out
}

func parseArg(name, typ string) *token.Arg {
	return &token.Arg{Name: name, Type: parseSig(typ)}
}

func parseSig(sig string) string {
	s, rlen := next(sig)
	if len(sig) != rlen {
		panic(fmt.Sprintf("warn: %q invalid signature", sig))
	}
	return s
}

func next(s string) (string, int) {
	if len(s) == 0 {
		return "", 0
	}
	switch s[0] {
	case 'y':
		return "byte", 1
	case 'b':
		return "bool", 1
	case 'n':
		return "int16", 1
	case 'q':
		return "uint16", 1
	case 'i':
		return "int32", 1
	case 'u':
		return "uint32", 1
	case 'x':
		return "int64", 1
	case 't':
		return "uint64", 1
	case 'd':
		return "float64", 1
	case 'h':
		return "dbus.UnixFD", 1
	case 's':
		return "string", 1
	case 'o':
		return "dbus.ObjectPath", 1
	case 'v':
		return "dbus.Variant", 1
	case 'g':
		return "dbus.Signature", 1
	case 'a':
		if s[1] == '{' { // dictionary
			i := 4
			k, rlen := next(s[2:])
			if rlen != 1 {
				panic("key is not a primitive")
			}
			v, rlen := next(s[3:])
			if rlen == 0 {
				panic("value is not available")
			}
			i += rlen
			return "map[" + k + "]" + v, i
		}
		s, rlen := next(s[1:])
		return "[]" + s, rlen + 1
	case '(':
		i := 1
		n := 1
		for i < len(s) && n != 0 {
			if s[i] == '(' {
				n++
			} else if s[i] == ')' {
				n--
			}
			i++
		}
		return "struct {" + strings.Join(structFields(s[1:i-1]), ";") + "}", i
	default:
		panic("not supported signature: " + string(s[0]))
	}
}

func structFields(sig string) []string {
	fields := make([]string, 0, len(sig))
	for i, v := 0, 0; i < len(sig); v++ {
		s, rlen := next(sig[i:])
		if rlen == 0 {
			break
		}
		i += rlen
		fields = append(fields, fmt.Sprintf("V%d %s", v, s))
	}
	return fields
}
