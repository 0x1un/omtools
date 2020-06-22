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
	"strings"
	"zbxtools"

	"github.com/buger/goterm"
	"github.com/spf13/cobra"
)

var (
	shells = map[string]func(string, string){
		"hostlist": lscmd,
	}
)

func shellcmd(cmd *cobra.Command, args []string) {
	buf := bufio.NewReader(os.Stdin)
	his := []string{}
	url, username, password := getInputWithPromptui()
	zbx = zbxtools.NewZbxTool(fmt.Sprintf("http://%s/api_jsonrpc.php", url), username, password)
	goterm.Clear()
	for i := 0; ; i++ {
		fmt.Printf("[%d]zbxtool~$ ", i)
		read, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		read = trimRightSapce(read)
		his = append(his, fmt.Sprintf("%d %s", i, read))
		lines := strings.Split(read, " ")[:]
		var cm string
		if len(lines) != 0 {
			cm = lines[0]
			if cm == "exit" {
				goterm.Clear()
				break
			}
			if cm == "his" {
				fmt.Println(strings.Join(his, "\n"))
				continue
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
