package cmd

import (
	"errors"
	"log"
	"omtools/adtools"
	"strings"
	"unicode"

	"github.com/go-ldap/ldap/v3"
	"github.com/manifoldco/promptui"
)

func findLongStringLength(s []string) int {
	l := 0
	for _, v := range s {
		if len(v) > l {
			l = len(v)
		}
	}
	return l
}

func findElement(s string, ss []string) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}

func trimRightSapce(str string) string {
	return strings.TrimRight(str, "\r\n")
}

func spaceStringsBuilder(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func getInputWithPromptui(a string) (string, string, string) {
	prompt := promptui.Prompt{
		Label:   "address",
		Default: "localhost",
	}
	url, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}
	prompt.Label = "Username"
	prompt.Default = "Admin"
	username, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	prompt.Label = "Password"
	prompt.Mask = '*'
	prompt.Default = "9x14fals"
	password, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	if a == "ad" {
		url = "ldap://" + url
	}
	return url, username, password
}

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
