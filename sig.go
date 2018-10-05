package dbusgen

import "strings"

func sigZeroValue(sig string) string {
	switch sig[0] {
	case 'b':
		return "false"
	case 'y', 'n', 'q', 'i', 'u', 'x', 't', 'd', 'h':
		return "0"
	case 's', 'o':
		return `""`
	case 'v', 'a':
		return "nil"
	case 'g':
		return "dbus.Signature{}"
	case '(':
		panic("not implemented yet") // TODO
	default:
		panic("not supported signature: " + string(sig[0]))
	}
}

type sig []string

func (s sig) join(sep string) string {
	return strings.Join([]string(s), sep)
}

func newSig(s string) sig {
	var ss []string
	for i := 0; i < len(s); {
		s, rlen := next(s[i:])
		if rlen == 0 {
			break
		}
		i += rlen
		ss = append(ss, s)
	}
	return ss
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
		return "interface{}", 1
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
		return "struct {" + newSig(s[1:i-1]).join(";") + "}", i
	default:
		panic("not supported signature: " + string(s[0]))
	}
}
