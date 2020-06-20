package zbxtools

import (
	"fmt"
	"io/ioutil"
	"log"
	"zbxtools/go-zabbix"
)

type ZbxTool struct {
	session *zabbix.Session
}

func NewZbxTool(url, username, password string) *ZbxTool {
	session, err := zabbix.NewSession(url, username, password)
	if err != nil {
		panic(err)
	}
	return &ZbxTool{
		session: session,
	}
}

// ExportAnyHosts export any hosts
func (z *ZbxTool) ExportAnyHosts(path, format string) error {
	hostParams := zabbix.HostGetParams{}
	hostParams.OutputFields = []string{"hostid"}

	hosts, err := z.session.GetHosts(hostParams)
	if err != nil {
		return fmt.Errorf("Error getting Hosts: %v", err)
	}

	if len(hosts) == 0 {
		log.Fatal("No Hosts found")
	}
	hostIDS := make([]string, 0)
	for _, host := range hosts {
		if len(host.HostID) == 0 {
			continue
		}
		hostIDS = append(hostIDS, host.HostID)
	}

	params := zabbix.ConfigurationParamsRequest{
		Format: format,
		Options: zabbix.ConfiguraOption{
			Hosts: hostIDS,
		},
	}
	respData, err := z.session.ConfiguraExport(params)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, []byte(respData), 0666)
	if err != nil {
		return err
	}
	return nil
}
