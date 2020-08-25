package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	g "omtools/zbxgraph"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

func main() {
	wg := sync.WaitGroup{}
	cfg, err := ini.Load("init.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	// create graph directory if it doesn't exist
	path := "./graph/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0777)
		if err != nil {
			logrus.Fatal(err)
		}
	}
	general := cfg.Section("GENERAL")
	graphs := cfg.Section("GRAPH")
	session := g.NewSeesion(
		general.Key("ZBX_HOST").String()+general.Key("ZBX_PORT").String(),
		general.Key("ZBX_USERNAME").String(),
		general.Key("ZBX_PASSWORD").String())
	err = session.Login()
	if err != nil {
		logrus.Fatal(err)
	}
	wg.Add(len(graphs.KeysHash()))
	for name, graphid := range graphs.KeysHash() {
		go func(id, name string) {
			data, err := session.DownloadTrafficGraph(id, general.Key("TIME_FROM").String(), general.Key("TIME_TO").String())
			if err != nil {
				logrus.Fatal(err)
			}
			err = ioutil.WriteFile(path+name+".png", data, 0644)
			if err != nil {
				logrus.Fatal(err)
			}
			wg.Done()
		}(graphid, name)
	}
	wg.Wait()
	fmt.Print("Press 'Enter' to continue...")
	_, err = bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		logrus.Fatal(err)
	}
}
