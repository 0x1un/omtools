package zbxtools

import (
	"io/ioutil"
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
	hosts, err := z.ListHostID("")
	if err != nil {
		return err
	}
	hostIDS := make([]string, 0)
	for hostid, _ := range hosts {
		hostIDS = append(hostIDS, hostid)
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
