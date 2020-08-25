package zbxgraph

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/ini.v1"
)

func TestLogin(t *testing.T) {
	cfg, err := ini.Load("testfiles/init.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	general := cfg.Section("GENERAL")
	graphs := cfg.Section("GRAPH")
	s := NewSeesion(general.Key("ZBX_HOST").String()+general.Key("ZBX_PORT").String(), general.Key("ZBX_USERNAME").String(), general.Key("ZBX_PASSWORD").String())
	err = s.Login()
	if err != nil {
		t.Fatal(err)
	}
	for name, graphid := range graphs.KeysHash() {
		data, err := s.DownloadTrafficGraph(graphid, "now-24h", "now")
		if err != nil {
			t.Fatal(err)
		}
		err = ioutil.WriteFile("graph/"+name+".png", data, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
}
