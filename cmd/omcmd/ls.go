package cmd

import (
	"github.com/spf13/cobra"
)

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
