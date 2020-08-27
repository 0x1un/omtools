package zbxtools

import (
	"fmt"
	"strings"

	"github.com/0x1un/omtools/go-zabbix"
)

func (z *ZbxTool) ListGroup(key string) (map[string]string, error) {

	params := zabbix.HostgroupGetParams{}

	hostgroups, err := z.session.GetHostgroups(params)
	if err != nil {
		return nil, fmt.Errorf("Error getting Hostgroups: %v", err)
	}

	if len(hostgroups) == 0 {
		return nil, fmt.Errorf("No Hostgroups found")
	}

	groupMap := make(map[string]string, 0)
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
