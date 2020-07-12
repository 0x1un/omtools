package cmd

import (
	"errors"
	"omtools/adtools"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-ldap/ldap/v3"
	"github.com/manifoldco/promptui"
)

func getUserInfo() (disName string, username string, org string, pwd string, descpt string, disabled bool) {
	validate := func(input string) error {
		if len(input) == 0 {
			return errors.New("你必须输入这个值!")
		}
		return nil
	}
	keys := []string{
		"Display Name",
		"Username",
		"Organization",
		"Password",
		"Description",
		"Disable",
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
