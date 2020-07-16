package cmd

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

func lscmd(key, target string) {
	tw := table.NewWriter()
	switch target {
	case "host":
		hosts, err := zbx.ListHostID(key)
		if err != nil {
			fmt.Println(err)
			return
		}
		tw.AppendHeader(table.Row{"id", "name", "ip"})
		rows := make([]table.Row, 0)
		for _, host := range hosts {
			rows = append(rows, table.Row{host.ID, host.Name, host.IP})
		}
		tw.AppendRows(rows)
		tw.SortBy([]table.SortBy{{Name: "id", Mode: table.AscNumeric}})
		fmt.Println(tw.Render())
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
		tw.AppendHeader(table.Row{"id", "name"})
		rows := make([]table.Row, 0)
		l := findLongStringLength(names)
		for id, name := range groupmap {
			lens := l - len(name)
			for i := 0; i < lens; i++ {
				name += " "
			}
			rows = append(rows, table.Row{id, name})
		}
		tw.AppendRows(rows)
		tw.SortBy([]table.SortBy{{Name: "id", Mode: table.AscNumeric}})
		fmt.Println(tw.Render())
	}
}

func zbxCmdHandler(line string) {
	switch {
	case strings.HasPrefix(line, "create host from "):
		if l := line[17:]; len(l) > 0 {
			err := zbx.CreateMultipleHost(l)
			if err != nil {
				fmt.Println(err)
			}
		}
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
