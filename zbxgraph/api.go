package zbxgraph

import (
	"fmt"
	"log"
	"strings"

	"github.com/0x1un/omtools/go-zabbix"
)

type zbxGraph struct {
	session *zabbix.Session
}

// 创建一个流量图获取的公用接口
func NewZbxGraph(url, username, password string) *zbxGraph {
	session, err := zabbix.NewSession(url, username, password)
	if err != nil {
		log.Panic(err)
	}
	return &zbxGraph{
		session: session,
	}
}

// 根据hostid 获取主机的所有流量图参数，包含有流量图
func (z *zbxGraph) GetGraphParameters(hostid []string) (zabbix.GraphResult, error) {
	params := zabbix.GraphGetParameters{}
	params.Output = "extend"
	params.Hostids = hostid
	params.SortField = []string{"name"}
	graphRes, err := z.session.GraphGet(params)
	if err != nil {
		return nil, err
	}
	return graphRes, nil
}

// return: group id :: group name
func (z *zbxGraph) ListGroup(key string) (map[string]string, error) {
	params := zabbix.HostgroupGetParams{}
	hostgroups, err := z.session.GetHostgroups(params)
	if err != nil {
		return nil, fmt.Errorf("Error getting Hostgroups: %v", err)
	}

	if len(hostgroups) == 0 {
		return nil, fmt.Errorf("No Hostgroups found")
	}
	groupMap := make(map[string]string)
	for _, hostgroup := range hostgroups {
		if hostgroup.GroupID == "" {
			continue
		}
		if len(key) != 0 {
			if strings.Contains(hostgroup.Name, key) {
				groupMap[hostgroup.GroupID] = hostgroup.Name
			}
		} else {
			groupMap[hostgroup.GroupID] = hostgroup.Name
		}

	}
	return groupMap, nil
}
