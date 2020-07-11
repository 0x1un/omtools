/*
Copyright © 2020 0x1un <aumujun@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"omtools/adtools"
	"omtools/zbxtools"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
)

const (
	zbxUrl = "http://%s/api_jsonrpc.php"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "交互模式",
	Long:  `进入交互模式`,
	Run:   shellcmd,
}

func init() {
	rootCmd.AddCommand(shellCmd)
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

func shellcmd(cmd *cobra.Command, args []string) {

	// set line prompt
	l, err := readline.NewEx(&readline.Config{
		Prompt:              "\033[31momtools »\033[0m ",
		HistoryFile:         "./readline.tmp",
		AutoComplete:        completer,
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	setPasswordCfg := l.GenPasswordConfig()
	setPasswordCfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		l.SetPrompt(fmt.Sprintf("Enter password(%v): ", len(line)))
		l.Refresh()
		return nil, 0, false
	})

	log.SetOutput(l.Stderr())
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		if line == "bye" || line == "exit" {
			goto exit
		}

		line = strings.TrimSpace(line)
		if line == "go zbx" {
		}
		switch line {
		case "go zbx":
			url, username, password := getInputWithPromptui("")
			zbx = zbxtools.NewZbxTool(fmt.Sprintf(zbxUrl, url), username, password)
		case "go ad":
			url, buser, bpass := getInputWithPromptui("ad")
			ad, err = adtools.NewADTools(url, buser, bpass)
			if err != nil {
				fmt.Printf("failed connect to %s, err:%s\n", url, err.Error())
			}
		}
		if zbx != nil {
			zbxCmdHandler(line)
		}
		if ad != nil {
			adCmdHandler(line)
		}
	}
exit:
}

func adCmdHandler(line string) {
	switch {
	case line == "add single user":
		disname, username, org, pwd, des, disabled := getUserInfo()
		err := ad.AddUser(disname, username, org, pwd, des, disabled)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func zbxCmdHandler(line string) {
	switch {
	case strings.HasPrefix(line, "list "):
		switch line[5:] {
		case "host":
			cmdMap[line[:4]]("", line[5:])
		case "group":
			cmdMap[line[:4]]("", line[5:])
		}
	// query [host] by [key]
	case strings.HasPrefix(line, "query "):
		subcmd := line[6:]
		subList := strings.Split(subcmd, " ")
		if len(subList) >= 3 {
			if subList[1] == "by" {
				cmdMap[line[:5]](subList[2], subList[0])
			}
		}
	case strings.HasPrefix(line, "cfg "):
		subcmd := line[4:]
		subList := strings.Split(subcmd, " ")
		if len(subList) >= 2 {
			if subList[0] == "export" {
				cmdMap[line[:3]](subList[1], "")
			}
		}
	case line == "go zbx":
		println("connect to zabbix server...")
	default:
		log.Println("you said:", strconv.Quote(line))
	}
}
