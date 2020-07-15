package cmd

import (
	"fmt"
	"strings"
)

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
		fmt.Println("connect to zabbix server...")
	case line == "re con zbx":
		fmt.Println("reconnect to zabbix server...")
	case len(line) == 0:
	default:
		fmt.Printf(cmdNotFound, line)
	}
}
