package cmd

import (
	"errors"
	"fmt"
	"omtools/adtools"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-ldap/ldap/v3"
	"github.com/manifoldco/promptui"
)

var (
	validate = func(input string) error {
		if len(input) == 0 {
			return errors.New("你必须输入这个值!")
		}
		return nil
	}
)

func unlockUser(user string) error {
	if user == "" {
		return errors.New("user cannot be empty")
	}
	res, err := ad.QueryUser(BaseDN, adtools.LockedAllUserFilter, ldap.ScopeWholeSubtree)
	if err != nil {
		return err
	}
	for _, v := range res.Entries {
		x := strings.Split(strings.ToLower(v.DN), ",")
		// 判断DN路径最顶部的值是否为user
		if len(x) > 0 {
			if x[0] == strings.ToLower("cn="+user) {
				err := ad.UnlockUser(v.DN)
				if err != nil {
					return err
				}
				fmt.Printf("unlock user: %s\n", v.DN)
				return nil
			}
		}
	}
	return nil
}

func changeStatus(user string, disabled bool) error {
	if user == "" {
		return errors.New("user cannot be empty")
	}
	res, err := ad.QueryUser(BaseDN, adtools.Ft(adtools.UserFilter, user, user), ldap.ScopeWholeSubtree)
	if err != nil {
		return err
	}
	for _, v := range res.Entries {
		x := strings.Split(strings.ToLower(v.DN), ",")
		// 判断DN路径最顶部的值是否为user
		if len(x) > 0 {
			if x[0] == strings.ToLower("cn="+user) {
				err := ad.ChangeUserStatus(v.DN, disabled)
				if err != nil {
					return err
				}
				fmt.Printf("change status: %s: %s", user, func(a bool) string {
					if a == false {
						return "Enabled"
					}
					return "Disabled"
				}(disabled))
				return nil
			}
		}
	}
	return nil
}

func getOuPath() string {
	p := promptui.Prompt{
		Label:    "OU key",
		Validate: validate,
	}
	res, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	queryRes, err := ad.QueryUser(BaseDN, adtools.Ft(adtools.OuWithoutDefaultOUFilter, res), ldap.ScopeWholeSubtree)
	if err != nil {
		log.Fatal(err)
	}

	list := []string{}

	for _, v := range queryRes.Entries {
		list = append(list, v.DN)
	}
	pSelect := promptui.Select{
		Label: "请选择目标",
		Items: list,
	}
	_, res, err = pSelect.Run()
	if err != nil {
		log.Fatal(err)
	}
	idx := strings.Index(strings.ToLower(res), strings.ToLower(BaseDN))
	res = res[:idx-1]
	return res
}

func getUserInfo() (disName string, username string, org string, pwd string, descpt string, disabled bool) {
	keys := []string{
		"Display Name", "Username", "Organization", "Password", "Description", "Disable",
	}
	info := []string{}
	for i := 0; i < 6; i++ {
		p := promptui.Prompt{
			Label: keys[i],
			Validate: func(a string) func(string) error {
				if a != "Disable" {
					return validate
				}
				return nil
			}(keys[i]),
			Mask: func(a string) rune {
				if a == "Password" {
					return '*'
				}
				return rune(0)
			}(keys[i]),
			IsConfirm: func(a string) bool {
				return a == "Disable"
			}(keys[i]),
		}
		res, err := p.Run()
		if err != nil && keys[i] == "Disable" {
			res = "n"
		} else if err != nil {
			log.Fatal(err)
		}
		if keys[i] == "Organization" {
			queryRes, err := ad.QueryUser(BaseDN, adtools.Ft(adtools.OuWithoutDefaultOUFilter, res), ldap.ScopeWholeSubtree)
			if err != nil {
				log.Fatal(err)
			}
			entrys := []string{}
			for _, v := range queryRes.Entries {
				entrys = append(entrys, v.DN)
			}
			pSelect := promptui.Select{
				Label: "Please select a target",
				Items: entrys,
			}
			_, res, err = pSelect.Run()
			if err != nil {
				log.Fatal(err)
			}
			idx := strings.Index(strings.ToLower(res), strings.ToLower(BaseDN))
			res = res[:idx-1]
		}
		info = append(info, res)
	}
	disName = info[0]
	username = info[1]
	org = info[2]
	pwd = info[3]
	descpt = info[4]
	disabled = info[5] == "y"
	return
}

func adCmdHandler(line string) {
	switch {
	case line == "add single user":
		disname, username, org, pwd, des, disabled := getUserInfo()
		err := ad.AddUser(disname, username, org, pwd, des, disabled)
		if err != nil {
			fmt.Println(err)
			return
		}
	case strings.HasPrefix(line, "add user from "):
		// TODO: 检查文件路径合法性
		l := line[14:]
		if len(l) == 0 {
			println("请输入文件路径")
			return
		}
		for _, e := range ad.AddUserMultiple(l, getOuPath(), false).Errors {
			fmt.Println(e)
		}
	case strings.HasPrefix(line, "del user with "):
		l := line[14:]
		if len(l) == 0 {

		}
	case strings.HasPrefix(line, "query info "):
		l := ""
		if l = line[11:]; l == "all" {
			l = "*"
		}
		res, err := ad.GetUserInfoTable(l)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(res)
	case strings.HasPrefix(line, "query user by "):
		// 以某种关键字筛选用户: 锁定状态、登入次数、激活状态

	case strings.HasPrefix(line, "dis ") || strings.HasPrefix(line, "ena "):
		l := line[4:]
		if len(l) != 0 {
			l = strings.TrimSpace(l)
			err := changeStatus(l, line[:3] == "dis")
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	case strings.HasPrefix(line, "unlock "):
		l := line[7:]
		if len(l) != 0 {
			l = strings.TrimSpace(l)
			if err := unlockUser(l); err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	case line == "go ad":
		fmt.Println("connect to ad server...")
	case line == "re con ad":
		fmt.Println("reconnect to ad server...")
	case len(line) == 0:
	default:
		fmt.Printf(cmdNotFound, line)
	}
}
