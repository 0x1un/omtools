/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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

	"github.com/spf13/cobra"
)

func lscmd(key, target string) {
	switch target {
	case "host":
		err := zbx.ListHostID(key)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "group":
		groupmap, err := zbx.ListGroup(key)
		if err != nil {
			fmt.Println(err)
			return
		}
		names := make([]string, len(groupmap))
		for _, name := range groupmap {
			names = append(names, name)
		}
		l := findLongStringLength(names)
		for id, name := range groupmap {
			lens := l - len(name)
			for i := 0; i < lens; i++ {
				name += " "
			}
			fmt.Printf("%s\t\t->%s\n", name, id)
		}
	}
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "列出相关的内容",
	Long:  `列出zabbix上一些数据，例如：主机，图形，模板等等`,
	Run: func(cmd *cobra.Command, args []string) {
		status, err := cmd.Flags().GetString("grep")
		if err != nil {
			panic(err)
		}
		target, err := cmd.Flags().GetString("target")
		if err != nil {
			panic(err)
		}
		lscmd(status, target)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().String("grep", "", "过滤出包含指定关键字的主机 (--grep sangfor)")
	lsCmd.Flags().String("target", "", "想要列出的目标 (--target host)")
}
