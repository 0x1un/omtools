package zbxtools

import (
	"fmt"
	"omtools/go-zabbix"
	"strings"
)

type hostinfo struct {
	ID   string
	IP   string
	Name string
}

func (z *ZbxTool) getAnyHostID() ([]string, error) {
	params := zabbix.HostGetParams{}
	params.OutputFields = []string{"hostid"}
	hosts, err := z.session.GetHosts(params)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for _, host := range hosts {
		if host.HostID == "" {
			continue
		}
		ids = append(ids, host.HostID)
	}
	return ids, nil
}

// @params key: 模糊搜索的关键字
func (z *ZbxTool) ListHostID(key string) ([]hostinfo, error) {
	hostParams := zabbix.HostGetParams{}
	hostParams.OutputFields = []string{"hostid", "name"}

	hosts, err := z.session.GetHosts(hostParams)
	if err != nil {
		return nil, fmt.Errorf("Error getting Hosts: %v", err)
	}

	if len(hosts) == 0 {
		return nil, fmt.Errorf("No Hosts found")
	}

	p2 := zabbix.HostInterfaceGetParams{}
	hostArray := make([]hostinfo, 0)
	for _, host := range hosts {
		if len(host.HostID) == 0 {
			continue
		}
		ip := func(id string) string {

			p2.Hostids = []string{host.HostID}
			p2.Output = []string{"ip"}
			ifres, err := z.session.HostInterfaceGet(p2)
			if err != nil {
				return ""
			}
			return ifres[0].IP
		}
		if len(key) == 0 {
			hostArray = append(hostArray, hostinfo{
				ID:   host.HostID,
				IP:   ip(host.HostID),
				Name: host.DisplayName,
			})
			continue
		}
		if strings.Contains(host.DisplayName, key) {
			hostArray = append(hostArray, hostinfo{
				ID:   host.HostID,
				IP:   ip(host.HostID),
				Name: host.DisplayName,
			})
		}
	}
	return hostArray, nil
}
