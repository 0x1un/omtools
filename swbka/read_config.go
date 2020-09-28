package main

import (
	"errors"
	"gopkg.in/ini.v1"
	"strings"
)

// readConfig 读取配置文件
func (s *swbka) readConfig(path string) (map[string]mulparam, error) {
	mp := make(map[string]mulparam)
	cfg, err := ini.LoadSources(
		ini.LoadOptions{SkipUnrecognizableLines: true,
			IgnoreInlineComment: true},
		path)
	if err != nil {
		return nil, err
	}
	s.defaultCFG.pubUser = cfg.Section("general").Key("pub_user").String()
	s.defaultCFG.pubPass = cfg.Section("general").Key("pub_pass").String()
	s.defaultCFG.pubTarget = strings.Split(cfg.Section("general").Key("pub_target").String(), ",")
	s.defaultCFG.projectName = cfg.Section("general").Key("project_name").String()
	// 如果没有指定拉取的配置文件，则默认拉取 "startup.cfg"
	if len(s.defaultCFG.pubTarget) == 0 {
		s.defaultCFG.pubTarget = []string{"startup.cfg"}
	}
	s.defaultCFG.pubPort = cfg.Section("general").Key("pub_port").String()
	s.defaultCFG.dingAtUsers = cfg.Section("general").Key("ding_at_users").Strings(",")
	s.defaultCFG.dingTokens = cfg.Section("general").Key("ding_tokens").Strings(",")
	s.defaultCFG.dingNotifyAll, _ = cfg.Section("general").Key("ding_notify_all").Bool()
	s.defaultCFG.profilePATH = cfg.Section("general").Key("profile_path").String()
	// webdav配置
	s.defaultCFG.webdavURL = cfg.Section("general").Key("webdav_url").String()
	s.defaultCFG.webdavUSER = cfg.Section("general").Key("webdav_user").String()
	s.defaultCFG.webdavPWD = cfg.Section("general").Key("webdav_pwd").String()
	// 邮件配置
	s.defaultCFG.smtpServer = cfg.Section("general").Key("smtp_server").String()
	s.defaultCFG.smtpPort, _ = cfg.Section("general").Key("smtp_port").Int()
	s.defaultCFG.smtpUSER = cfg.Section("general").Key("smtp_user").String()
	s.defaultCFG.smtpPWD = cfg.Section("general").Key("smtp_pwd").String()
	s.defaultCFG.smtpFROM = cfg.Section("general").Key("smtp_from").String()
	s.defaultCFG.smtpTO = cfg.Section("general").Key("smtp_to").Strings(",")
	for _, v := range cfg.Sections() {
		name := v.Name()
		if name == "general" || name == "DEFAULT" {
			continue
		}
		m := mulparam{}
		for ip, loginStr := range v.KeysHash() {
			s.total++
			strList := strings.Split(loginStr, ",")
			p := param{}
			p.ip = ip + s.defaultCFG.pubPort
			// 当未设置任何参数时，全部使用general中的配置
			if len(strList) == 0 {
				p.password = s.defaultCFG.pubPass
				p.target = s.defaultCFG.pubTarget
				p.deviceName = "unknown"
				p.username = s.defaultCFG.pubUser
			}
			// 当其中任何一个参数被设置时，则应用该参数，并使用general中的配置填补空缺的参数
			for idx, arg := range strList {
				if strings.TrimSpace(arg) == "" {
					switch idx {
					case 0:
						p.username = s.defaultCFG.pubUser
					case 1:
						p.password = s.defaultCFG.pubPass
					case 2:
						p.target = s.defaultCFG.pubTarget
					case 3:
						p.deviceName = "unknown"
					}
				} else {
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
			if len(p.target) == 0 ||
				p.password == "" || p.username == "" {
				return nil, errors.New("missing matched the general config")
			}
			m.profiles = append(m.profiles, p)
		}
		mp[name] = m
	}
	return mp, nil
}
