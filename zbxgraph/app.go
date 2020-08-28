package zbxgraph

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

func Run(config interface{}, prefix string, randName bool) (map[string][]string, error) {
	wg := sync.WaitGroup{}
	cfg, err := ini.Load(config)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	// create graph directory if it doesn't exist
	graphPrefix := prefix
	if _, err := os.Stat(graphPrefix); os.IsNotExist(err) {
		err := os.Mkdir(graphPrefix, 0777)
		if err != nil {
			return nil, err
		}
	}
	general := cfg.Section("GENERAL")
	session := NewSeesion(
		general.Key("ZBX_HOST").String()+general.Key("ZBX_PORT").String(),
		general.Key("ZBX_USERNAME").String(),
		general.Key("ZBX_PASSWORD").String())
	err = session.Login()
	if err != nil {
		return nil, err
	}
	outputMap := map[string][]string{}
	for _, section := range cfg.Sections() {
		if section.Name() == "Default" || len(section.KeysHash()) == 0 ||
			section.Name() == "GENERAL" || section.Name() == "INSPECTION" {
			continue
		}
		graphLocalPath := graphPrefix + section.Name()
		if _, err := os.Stat(graphLocalPath); os.IsNotExist(err) {
			err := os.Mkdir(graphLocalPath, 0777)
			if err != nil {
				return nil, err
			}
		}
		wg.Add(len(section.KeysHash()))
		for name, graphid := range section.KeysHash() {
			go func(n, g string) {
				data, err := session.DownloadTrafficGraph(g, general.Key("TIME_FROM").String(), general.Key("TIME_TO").String())
				if err != nil {
					logrus.Println(err)
				}
				rn := ""
				if randName {
					rn = fmt.Sprintf("%d", rand.Int())
				}
				outputFile := graphLocalPath + "/" + n + "-" + rn + ".png"
				err = ioutil.WriteFile(outputFile, data, 0644)
				if err != nil {
					logrus.Println(err)
				}
				outputMap[n] = append(outputMap[n], outputFile)
				wg.Done()
			}(name, graphid)
		}
		wg.Wait()
	}

	return outputMap, nil
}
