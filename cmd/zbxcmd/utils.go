package cmd

import (
	"strings"
	"unicode"

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

func getInputWithPromptui() (string, string, string) {
	prompt := promptui.Prompt{
		Label:   "host",
		Default: "localhost",
	}
	url, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	prompt.Label = "Username"
	prompt.Default = "Admin"
	username, err := prompt.Run()
	if err != nil {
		panic(err)
	}

	prompt.Label = "Password"
	prompt.Mask = '*'
	prompt.Default = "zabbix"
	password, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	return url, username, password
}
