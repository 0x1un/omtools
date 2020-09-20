package cmd

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"
)

func expAnyHosts(filename, _ string) {

	// 获取文件的后缀
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
	_ = args
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
