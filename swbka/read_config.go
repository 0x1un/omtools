package main

import (
	"gopkg.in/ini.v1"
	"strings"
	"errors"
)

// readConfig 读取配置文件
func (*swbka) readConfig(path string) (map[string]mulparam, error) {
	mp := make(map[string]mulparam)
	cfg, err := ini.LoadSources(
		ini.LoadOptions{SkipUnrecognizableLines: true,
			IgnoreInlineComment: true},
		path)
	if err != nil {
		return nil, err
	}
	pubUser := cfg.Section("general").Key("pub_user").String()
	pubPass := cfg.Section("general").Key("pub_pass").String()
	pubTarget := cfg.Section("general").Key("pub_target").String()
	pubPort := cfg.Section("general").Key("pub_port").String()
	// 如果没有指定拉取的配置文件，则默认拉取 "startup.cfg"
	if pubTarget == "" {
		pubTarget = "startup.cfg"
	}
	for _, v := range cfg.Sections() {
		name := v.Name()
		if name == "general" || name == "DEFAULT" {
			continue
		}
		m := mulparam{}
		for ip, loginStr := range v.KeysHash() {
			strList := strings.Split(loginStr, ",")
			p := param{}
			p.ip = ip+pubPort
			// 当未设置任何参数时，全部使用general中的配置
			if len(strList) == 0 {
				p.password = pubPass
				p.target = strings.Split(pubTarget, ",")
				p.deviceName = "unknown"
				p.username = pubUser
			}
			// 当其中任何一个参数被设置时，则应用该参数，并使用general中的配置填补空缺的参数
			for idx, arg := range strList {
				if strings.TrimSpace(arg) == "" {
					switch idx {
					case 0:
						p.username = pubUser
					case 1:
						p.password = pubPass
					case 2:
						p.target = strings.Split(pubTarget, ",")
					case 3:
						p.deviceName = "unknown"
					}
				}else {
					switch idx {
					case 0:
						p.username = strList[0]
					case 1:
						p.password = strList[1]
					case 2:
						p.target = []string{strList[2]}
					case 3:
						p.deviceName = strList[3]
					}
				}
			}
			if  len(p.target) == 0 ||
				p.password == "" || p.username == "" {
				return nil, errors.New("missing matched the general config")
			}
			m.profiles = append(m.profiles, p)
		}
		mp[name] = m
	}
	return mp, nil
}
