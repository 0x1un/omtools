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
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"zbxtools"

	"github.com/buger/goterm"
	"github.com/spf13/cobra"
)

var (
	shells = map[string]func(string, string){
		"hostlist": lscmd,
	}
	commands = []string{"host", "exit", "clear", "clean", "his", "group"}
	help     = map[string]string{
		"host": `host query [xxx]
host list all`,
		"group": `group query [xxx]
group list all`,
	}
)

func shellcmd(cmd *cobra.Command, args []string) {
	buf := bufio.NewReader(os.Stdin)
	his := []string{}
	url, username, password := getInputWithPromptui()
	zbx = zbxtools.NewZbxTool(fmt.Sprintf("http://%s/api_jsonrpc.php", url), username, password)
loop:
	for i := 0; ; i++ {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for sig := range c {
				// handle it
			}
		}()
		fmt.Printf("[%d]zbxtool~$ ", i)
		read, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		read = trimRightSapce(read)
		his = append(his, fmt.Sprintf("%d %s", i, read))
		lines := strings.Split(read, " ")[:]
		var cm string
		if l := len(lines); l != 0 {
			cm = lines[0]
			if !findElement(cm, commands) && !(len(cm) == 0) {
				fmt.Printf("Unkown command: %s\n", cm)
				continue
			}
			switch cm {
			case "exit":
				println("bye~~")
				break loop
			case "his":
				println(strings.Join(his, "\n"))
				continue
			case "clear", "clean":
				goterm.MoveCursor(1, 1)
				goterm.Clear()
			default:
				if l == 1 {
					println(help[cm])
				}
			}

			for i := 1; i < len(lines); i++ {
				switch lines[i] {
				case "list":
					if x := i + 1; x < len(lines) && lines[x] == "all" {
						lscmd("", cm)

					}
				case "query":
					if x := i + 1; x < len(lines) {
						lscmd(lines[x], cm)
					}
				}
			}
		}
	}
}

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
