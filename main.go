package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/amenzhinsky/dbus-codegen-go/parser"
	"github.com/amenzhinsky/dbus-codegen-go/printer"
	"github.com/amenzhinsky/dbus-codegen-go/token"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

var (
	destFlag       []string
	onlyFlag       []string
	exceptFlag     []string
	prefixFlag     []string
	systemFlag     bool
	packageFlag    string
	gofmtFlag      bool
	xmlFlag        bool
	outputFlag     string
	serverOnlyFlag bool
	clientOnlyFlag bool
	camelizeFlag   bool
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: dbus-codegen-go [option...] PATH...

D-Bus Introspection Data Format code generator for Golang.

Options:
`)
		flag.PrintDefaults()
	}
	flag.Var((*stringsVar)(&destFlag), "dest", "`destination` name(s) to introspect")
	flag.Var((*stringsVar)(&onlyFlag), "only", "generate code only for the named `interface`s")
	flag.Var((*stringsVar)(&exceptFlag), "except", "skip the named `interface`s")
	flag.Var((*stringsVar)(&prefixFlag), "prefix", "`prefix` to strip from interface names")
	flag.BoolVar(&systemFlag, "system", false, "connect to the system bus")
	flag.StringVar(&packageFlag, "package", "dbusgen", "generated package `name`")
	flag.BoolVar(&gofmtFlag, "gofmt", true, "gofmt results")
	flag.BoolVar(&xmlFlag, "xml", false, "combine the dest's introspections into a single document")
	flag.StringVar(&outputFlag, "output", "", "`path` to output destination")
	flag.BoolVar(&serverOnlyFlag, "server-only", false, "generate only server-side code")
	flag.BoolVar(&clientOnlyFlag, "client-only", false, "generate only client-side code")
	flag.BoolVar(&camelizeFlag, "camelize", false, "camelize type names omitting underscores")
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(destFlag) == 0 && xmlFlag {
		return errors.New("cannot combine -xml and -dest")
	}
	if serverOnlyFlag && clientOnlyFlag {
		return errors.New("cannot combine -server-only and -client-only")
	}
	if len(onlyFlag) != 0 && len(exceptFlag) != 0 {
		return errors.New("cannot combine -only and -except")
	}

	var ifaces []*token.Interface
	switch {
	case len(destFlag) != 0:
		if flag.NArg() > 0 {
			return errors.New("cannot combine -dest and file paths")
		}
		conn, err := connect(systemFlag)
		if err != nil {
			return err
		}
		defer conn.Close()

		if xmlFlag {
			b, err := generateXML(conn, destFlag)
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
	case flag.NArg() > 0:
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
	default:
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		ifaces, err = parser.Parse(b)
		if err != nil {
			return err
		}
	}

	filtered := make([]*token.Interface, 0, len(ifaces))
	for _, iface := range ifaces {
		if isNeeded(iface.Name) {
			filtered = append(filtered, iface)
		}
	}

	buf := &bytes.Buffer{}
	if err := printer.Print(buf, filtered,
		printer.WithPackageName(packageFlag),
		printer.WithGofmt(gofmtFlag),
		printer.WithPrefixes(prefixFlag),
		printer.WithServerOnly(serverOnlyFlag),
		printer.WithClientOnly(clientOnlyFlag),
		printer.WithCamelize(camelizeFlag),
	); err != nil {
		return err
	}

	output := os.Stdout
	if outputFlag != "" {
		f, err := os.OpenFile(outputFlag, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		output = f
	}
	_, err := io.Copy(output, buf)
	return err
}

func connect(system bool) (*dbus.Conn, error) {
	if system {
		return dbus.SystemBus()
	}
	return dbus.SessionBus()
}

func parseDest(conn *dbus.Conn, dests []string) ([]*token.Interface, error) {
	ifaces := make([]*token.Interface, 0, 16)
	for _, dest := range dests {
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
	}
	return ifaces, nil
}

func generateXML(conn *dbus.Conn, dests []string) ([]byte, error) {
	var ifaces []introspect.Interface
	for _, dest := range dests {
		if err := introspectDest(conn, dest, "/", func(n *introspect.Node) error {
			for _, ifn := range n.Interfaces {
				var found bool
				for _, ifc := range ifaces {
					if ifc.Name == ifn.Name {
						found = true
						break
					}
				}
				if !found && isNeeded(ifn.Name) {
					ifaces = append(ifaces, ifn)
				}
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}
	return xml.MarshalIndent(&introspect.Node{
		Interfaces: ifaces,
	}, "", "\t")
}

func merge(curr, next []*token.Interface) []*token.Interface {
	for _, ifn := range next {
		var found bool
		for _, ifc := range curr {
			if ifc.Name == ifn.Name {
				found = true
				break
			}
		}
		if !found {
			curr = append(curr, ifn)
		}
	}
	return curr
}

func isNeeded(iface string) bool {
	return len(onlyFlag) == 0 && len(exceptFlag) == 0 ||
		len(onlyFlag) != 0 && includes(onlyFlag, iface) ||
		len(exceptFlag) != 0 && !includes(exceptFlag, iface)
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
	node, err := introspect.Call(conn.Object(dest, path))
	if err != nil {
		return err
	}
	if err := fn(node); err != nil {
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

type stringsVar []string

func (ss *stringsVar) String() string {
	return "[" + strings.Join(*ss, ", ") + "]"
}

func (ss *stringsVar) Set(arg string) error {
	for _, s := range strings.Split(arg, ",") {
		if s == "" {
			continue
		}
		*ss = append(*ss, s)
	}
	return nil
}
