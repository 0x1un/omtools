package cmd

import (
	"fmt"
	"io"
	"omtools/adtools"
	"omtools/zbxtools"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
)

const (
	zbxUrl      = "http://%s/api_jsonrpc.php"
	cmdNotFound = "omtools: command not found: %s\n"
)

var (
	mode        = ""
	sessionInfo = map[string]string{}
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "交互模式",
	Long:  `进入交互模式`,
	Run:   shellcmd,
}

var destory = func() {
	if ad != nil {
		ad.BuiltinConn().Close()
		ad = nil
	}
}

func init() {
	rootCmd.AddCommand(shellCmd)
}

func shellcmd(cmd *cobra.Command, args []string) {

	// set line prompt
	l, err := readline.NewEx(&readline.Config{
		Prompt:              "\033[31momtools »\033[0m ",
		HistoryFile:         "./readline.tmp",
		AutoComplete:        completer,
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	setPasswordCfg := l.GenPasswordConfig()
	setPasswordCfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		l.SetPrompt(fmt.Sprintf("Enter password(%v): ", len(line)))
		l.Refresh()
		return nil, 0, false
	})

	log.SetOutput(l.Stderr())
	log.SetReportCaller(true)
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		if line == "bye" || line == "exit" {
			goto exit
		}

		line = strings.TrimSpace(line)
		switch line {
		case "help":
			usage(l.Stderr())
		case "re con zbx":
			if mode == "zbx" && zbx != nil && len(sessionInfo) != 0 {
				zbx = zbxtools.NewZbxTool(fmt.Sprintf(zbxUrl, sessionInfo["zbxAddr"]), sessionInfo["zbxUser"], sessionInfo["zbxPwd"])
			}

		case "re con ad":
			if mode == "ad" && ad != nil && len(sessionInfo) != 0 {
				ad, err = adtools.NewADTools(sessionInfo["adAddr"], sessionInfo["adUser"], sessionInfo["adPwd"])
				if err != nil {
					fmt.Printf("failed connect to %s, err:%s\n", sessionInfo["adAddr"], err.Error())
				}
			}
		case "go zbx":
			url, username, password := getInputWithPromptui("")
			zbx = zbxtools.NewZbxTool(fmt.Sprintf(zbxUrl, url), username, password)
			mode = "zbx"
			sessionInfo["zbxAddr"] = url
			sessionInfo["zbxUser"] = username
			sessionInfo["zbxPwd"] = password
		case "go ad":
			url, buser, bpass := getInputWithPromptui("ad")
			ad, err = adtools.NewADTools(url, buser, bpass)
			if err != nil {
				fmt.Printf("failed connect to %s, err:%s\n", url, err.Error())
			}
			mode = "ad"
			sessionInfo["adAddr"] = url
			sessionInfo["adUser"] = buser
			sessionInfo["adPwd"] = bpass
		}
		if mode == "zbx" && zbx != nil {
			zbxCmdHandler(line)
		}
		if mode == "ad" && ad != nil {
			adCmdHandler(line)

		}
	}
exit:
	destory()
}
