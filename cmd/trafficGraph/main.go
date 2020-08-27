package main

import (
	g "omtools/zbxgraph"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := g.Run("init.ini", "./graph/"); err != nil {
		logrus.Fatal(err)
	}
}
