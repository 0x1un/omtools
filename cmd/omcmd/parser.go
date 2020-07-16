package cmd

import (
	"io"
	"io/ioutil"

	"github.com/chzyer/readline"
)

var cmdMap = map[string]func(string, string){
	// zabbix commands
	"list":  lscmd,
	"query": lscmd,
	"cfg":   expAnyHosts,
}

func usage(w io.Writer) {
	_, _ = io.WriteString(w, "commands:\n")
	_, _ = io.WriteString(w, completer.Tree("    "))
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

// Function constructor - constructs new function for listing given directory
func listFiles(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("query",
		readline.PcItem("host", readline.PcItem("by")),
		readline.PcItem("tpl", readline.PcItem("by")),
		readline.PcItem("graph", readline.PcItem("by")),
		readline.PcItem("info", readline.PcItem("*"), readline.PcItem("all")),
		readline.PcItem("user", readline.PcItem("by")),
	),
	readline.PcItem("dis"),
	readline.PcItem("ena"),
	readline.PcItem("unlock"),
	readline.PcItem("re",
		readline.PcItem("con",
			readline.PcItem("ad"),
			readline.PcItem("zbx"))),
	readline.PcItem("add",
		readline.PcItem("single", readline.PcItem("user")),
		readline.PcItem("user", readline.PcItem("from")),
	),
	readline.PcItem("del",
		readline.PcItem("user", readline.PcItem("with")),
	),
	readline.PcItem("go",
		readline.PcItem("zbx"),
		readline.PcItem("ad")),
	readline.PcItem("list",
		readline.PcItem("host"),
	),
	readline.PcItem("login"),
	readline.PcItem("bye"),
	readline.PcItem("help"),
)
