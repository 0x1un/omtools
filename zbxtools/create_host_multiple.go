package zbxtools

import (
	"encoding/csv"
	"fmt"
	"log"
	"omtools/go-zabbix"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func (z *ZbxTool) CreateMultipleHost(filename string /*这是一个csv文件*/) error {
	reqs, err := parseCsv(filename)
	if err != nil {
		return err
	}
	for _, req := range reqs {
		res, err := z.session.CreateHost(req)
		if err != nil {
			return err
		}
		fmt.Println(res)
	}
	return nil
}

func parseCsv(filename string) ([]zabbix.CreateHostRequest, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, errors.Errorf("failed read csv file: %s, error: %s\n", filename, err.Error())
	}
	// 主机名称,显示名称,所属群组id,接口IP,模板ID,禁用主机,使用IP,DNS,描述
	// 接口IP = host:port+type+main
	req := []zabbix.CreateHostRequest{}
	for _, record := range records[1:] {
		interfaces := []zabbix.Interface{}
		ipss := strings.Split(strings.TrimSpace(record[3]), "&")
		for _, addrs := range ipss {
			iface := zabbix.Interface{}
			addr := strings.Split(addrs, "+")
			iface.Bulk = 1
			iface.IP = strings.Split(addr[0], ":")[0]
			iface.Port = parseInt(strings.Split(addr[0], ":")[1])
			// 如果在ip中找到了main标识，将其赋值，否则默认将此接口设为主要
			iface.Main = func(a []string) int {
				if len(addr) > 2 {
					return parseInt(addr[2])
				}
				return 1
			}(addr)
			iface.Type = func(a []string) int {
				if len(a) >= 2 {
					return parseInt(a[1])
				}
				return 1
			}(addr) // default 1 = agent interface
			iface.Useip = 1
			if record[6] == "0" {
				iface.Useip = 0
				iface.DNS = record[7] // if useip == 0, dns cannot be empty
			}
			interfaces = append(interfaces, iface)
		}
		req = append(req, zabbix.CreateHostRequest{
			Status:      parseInt(record[5]),
			Host:        record[0],
			VisibleName: record[1],
			Description: record[len(record)-1],
			Interfaces:  interfaces,
			Groups: func(a string) []zabbix.Group {
				aa := strings.Split(a, "+")
				as := []zabbix.Group{}
				for _, v := range aa {
					as = append(as, zabbix.Group{
						GroupID: v,
					})
				}
				return as
			}(record[2]),
		})
	}
	return req, nil
}

func parseInt(a string) int {
	res, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return int(res)
}
