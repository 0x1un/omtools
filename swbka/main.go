// 此程序不再维护

package main

import "github.com/sirupsen/logrus"

func main() {
	if err := Impl(); err != nil {
		logrus.Fatal(err)
	}
}
