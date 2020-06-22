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
	"unicode"
	"zbxtools"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	shells = map[string]func(string, string){
		"hostlist": lscmd,
	}
)

func spaceStringsBuilder(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "交互模式",
	Long:  `进入交互模式`,
	Run: func(cmd *cobra.Command, args []string) {
		buf := bufio.NewReader(os.Stdin)
		prompt := promptui.Prompt{
			Label:   "host",
			Default: "localhost",
		}
		url, err := prompt.Run()
		if err != nil {
			panic(err)
		}
		prompt.Label = "Username"
		prompt.Default = "Admin"
		username, err := prompt.Run()
		if err != nil {
			panic(err)
		}

		prompt.Label = "Password"
		prompt.Mask = '*'
		prompt.Default = "zabbix"
		password, err := prompt.Run()
		if err != nil {
			panic(err)
		}
		zbx = zbxtools.NewZbxTool(fmt.Sprintf("http://%s/api_jsonrpc.php", url), username, password)
		for {
			print("zbxtool~$ ")
			read, err := buf.ReadString('\n')
			if err != nil {
				break
			}
			line := spaceStringsBuilder(read)
			if line == "quit" || line == "exit" {
				break
			}
			if f, ok := shells[line]; ok {
				// grep := "|grep"
				f("VI", "host")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
