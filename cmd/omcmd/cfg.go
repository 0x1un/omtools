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
	"path"

	"github.com/spf13/cobra"
)

func expAnyHosts(filename, _ string) {

	fileExt := path.Ext(filename)
	err := zbx.ExportAnyHosts(filename, func(s string) string {
		switch s {
		case ".xml":
			return "xml"
		default:
			return "json"
		}
	}(fileExt))
	if err != nil {
		println(err)
	}
	fmt.Printf("exported %s done\n", filename)
}

func cfgcmd(cmd *cobra.Command, args []string) {

	argName, _ := cmd.Flags().GetString("export")
	if argName == "" {
		if err := cmd.Help(); err != nil {
			panic(err)
		}
		return
	}
	expAnyHosts(argName, "")
}

// cfgCmd represents the cfg command
var cfgCmd = &cobra.Command{
	Use:   "cfg",
	Short: "zbx配置导入导出",
	Long:  `你可以使用此命令导入或者导出zabbix配置，诸如：流量图、主机、主机群组等等`,
	Run:   cfgcmd,
}

func init() {
	rootCmd.AddCommand(cfgCmd)
	cfgCmd.Flags().String("export", "", "export zabbix config")
}
