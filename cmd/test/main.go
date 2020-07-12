package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
)

// const (
// 	LDAP     = "ldap://172.19.2.10"
// 	username = "administrator@0x1un.io"
// 	password = "gdlk@123"
// )

// var (
// 	conn adtools.ADTooller
// )

// func init() {
// 	con, err := adtools.NewADTools(LDAP, username, password)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	conn = con
// }

func main() {

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
		info = append(info, res)
	}
	fmt.Println(info)
}
