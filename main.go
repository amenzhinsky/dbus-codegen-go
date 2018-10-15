package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/amenzhinsky/dbus-codegen-go/parser"
	"github.com/amenzhinsky/dbus-codegen-go/printer"
	"github.com/amenzhinsky/dbus-codegen-go/token"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

var (
	destFlag     string
	onlyFlag     []string
	exceptFlag   []string
	prefixesFlag []string
	sessionFlag  bool
	packageFlag  string
	gofmtFlag    bool
	xmlFlag      bool
)

type stringsFlag []string

func (ss *stringsFlag) String() string {
	return "[" + strings.Join(*ss, ", ") + "]"
}

func (ss *stringsFlag) Set(s string) error {
	if s = strings.Trim(s, " "); s == "" {
		return errors.New("string is empty")
	}
	*ss = append(*ss, s)
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s [FLAG...] [PATH...]

Take D-Bus Introspection Data Format and generates go code for it.

Flags:
`, os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&destFlag, "dest", "", "DBus destination name to introspect")
	flag.Var((*stringsFlag)(&onlyFlag), "only", "generate code only for the named interfaces")
	flag.Var((*stringsFlag)(&exceptFlag), "except", "skip the named interfaces")
	flag.Var((*stringsFlag)(&prefixesFlag), "prefix", "prefix to strip from interface names")
	flag.BoolVar(&sessionFlag, "session", false, "connect to the session bus instead of the system")
	flag.StringVar(&packageFlag, "package", "dbusgen", "generated package name")
	flag.BoolVar(&gofmtFlag, "gofmt", true, "gofmt results")
	flag.BoolVar(&xmlFlag, "xml", false, "combine the dest's introspections into a single document")
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	var ifaces []*token.Interface
	if destFlag == "" && xmlFlag {
		return errors.New("flag -xml cannot be used without -dest flag")
	}
	if destFlag != "" {
		if flag.NArg() > 0 {
			return errors.New("cannot combine arguments and -dest flag")
		}
		conn, err := connect(sessionFlag)
		if err != nil {
			return err
		}
		defer conn.Close()

		if xmlFlag {
			b, err := generateXml(conn, destFlag)
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil
		}
		ifaces, err = parseDest(conn, destFlag)
		if err != nil {
			return err
		}
	} else if flag.NArg() > 0 {
		for _, filename := range flag.Args() {
			b, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}
			chunk, err := parser.Parse(b)
			if err != nil {
				return err
			}
			ifaces = merge(ifaces, chunk)
		}
	} else {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		ifaces, err = parser.Parse(b)
		if err != nil {
			return err
		}
	}

	if len(onlyFlag) != 0 && len(exceptFlag) != 0 {
		return errors.New("cannot combine -only and -except flags")
	}
	filtered := make([]*token.Interface, 0, len(ifaces))
	for _, iface := range ifaces {
		if len(onlyFlag) == 0 && len(exceptFlag) == 0 ||
			len(onlyFlag) != 0 && includes(onlyFlag, iface.Name) ||
			len(exceptFlag) != 0 && !includes(exceptFlag, iface.Name) {
			filtered = append(filtered, iface)
		}
	}
	return printer.Print(os.Stdout, filtered,
		printer.WithPackageName(packageFlag),
		printer.WithGofmt(gofmtFlag),
		printer.WithPrefixes(prefixesFlag),
	)
}

func connect(session bool) (*dbus.Conn, error) {
	if session {
		return dbus.SessionBus()
	}
	return dbus.SystemBus()
}

func parseDest(conn *dbus.Conn, dest string) ([]*token.Interface, error) {
	ifaces := make([]*token.Interface, 0, 16)
	if err := introspectDest(conn, dest, "/", func(node *introspect.Node) error {
		chunk, err := parser.ParseNode(node)
		if err != nil {
			return err
		}
		ifaces = merge(ifaces, chunk)
		return nil
	}); err != nil {
		return nil, err
	}
	return ifaces, nil
}

func generateXml(conn *dbus.Conn, dest string) ([]byte, error) {
	var node introspect.Node
	if err := introspectDest(conn, dest, "/", func(n *introspect.Node) error {
		for _, iface := range n.Interfaces {
			var found bool
			for _, ifc := range node.Interfaces {
				if ifc.Name == iface.Name {
					found = true
					break
				}
			}
			if !found {
				node.Interfaces = append(node.Interfaces, iface)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	b, err := xml.MarshalIndent(&node, "", "\t")
	if err != nil {
		return nil, err
	}
	return b, nil
}

func merge(curr, next []*token.Interface) []*token.Interface {
	for j := range next {
		var found bool
		for i := range curr {
			if curr[i].Name == next[j].Name {
				found = true
				break
			}
		}
		if !found {
			curr = append(curr, next[j])
		}
	}
	return curr
}

func includes(ss []string, s string) bool {
	for i := range ss {
		if ss[i] == s {
			return true
		}
	}
	return false
}

func introspectDest(
	conn *dbus.Conn, dest string, path dbus.ObjectPath,
	fn func(node *introspect.Node) error,
) error {
	var s string
	if err := conn.Object(dest, path).Call(
		"org.freedesktop.DBus.Introspectable.Introspect", 0,
	).Store(&s); err != nil {
		return err
	}
	var node introspect.Node
	if err := xml.Unmarshal([]byte(s), &node); err != nil {
		return err
	}
	if err := fn(&node); err != nil {
		return err
	}
	if path == "/" {
		path = ""
	}
	for _, child := range node.Children {
		if err := introspectDest(conn, dest, path+"/"+dbus.ObjectPath(child.Name), fn); err != nil {
			return err
		}
	}
	return nil
}
