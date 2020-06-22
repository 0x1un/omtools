package zbxtools

import (
	"fmt"
	"strings"
	"zbxtools/go-zabbix"
)

// @params key: 模糊搜索的关键字
func (z *ZbxTool) ListHostID(key string) (map[string]string, error) {
	hostParams := zabbix.HostGetParams{}
	hostParams.OutputFields = []string{"hostid", "name"}

	hosts, err := z.session.GetHosts(hostParams)
	if err != nil {
		return nil, fmt.Errorf("Error getting Hosts: %v", err)
	}

	if len(hosts) == 0 {
		return nil, fmt.Errorf("No Hosts found")
	}
	hostIDS := make(map[string]string, len(hosts))
	for _, host := range hosts {
		if len(host.HostID) == 0 {
			continue
		}
		if len(key) == 0 {
			hostIDS[host.HostID] = host.DisplayName
			continue
		}
		if strings.Contains(host.DisplayName, key) {
			hostIDS[host.HostID] = host.DisplayName
		}
	}
	return hostIDS, nil
}
