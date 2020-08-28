package main

import (
	"bufio"
	"fmt"
	"os"

	g "github.com/0x1un/omtools/zbxgraph"

	"github.com/sirupsen/logrus"
)

func main() {
	if outFile, err := g.Run("init.ini", "graph/", true); err != nil {
		logrus.Fatal(err)
	} else {
		fmt.Println(outFile)
	}
	fmt.Print("Press 'Enter' to continue...")
	_, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		logrus.Fatal(err)
	}
}
