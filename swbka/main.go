/*
这个程序的目的是为了备份交换机的配置文件，使用ftp来下载交换机上的配置文件。
Create date: 2020-09-19
Author: 0x1un
*/
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/0x1un/omtools/swbka/glimit"
	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
	"github.com/studio-b12/gowebdav"
)

// 通用变量初始化
var (
	timeFormat = "20060102_15_04"
	webdavURI  = flag.String("webdavuri", "http://10.100.100.242/dav", "webdav uri")
	webdavUSER = flag.String("webdavuser", "swbka@public.com", "webdav username")
	webdavPWD  = flag.String("webdavpass", "G9apckB4rcApHxEEhIr72dKJg7jd1PPf", "webdav password")
	configPath = flag.String("c", "/etc/swbka/profile.ini", "config file location")
)

// 错误类型初始化
var (
// diaWebdavErr
)

type param struct {
	ip         string
	username   string
	password   string
	target     []string
	deviceName string
}

type mulparam struct {
	profiles []param
}

// swbka: 入口结构
type swbka struct {
	sumFile string
	sumMap  sync.Map
	wd      *gowebdav.Client
}

func (s *swbka) pushConfig2Webdav(buf []byte, path string) error {
	return s.wd.Write(path, buf, 0644)
}


// downloadFileMock 模拟下载 测试使用
func (s *swbka) downloadFileMock(p param) error {
	fmt.Println(p.ip)
	time.Sleep(2 * time.Second)
	return nil
}

// downloadSwitchCfg 批量下载交换机配置文件
func (s *swbka) downloadSwitchCfg(sws map[string]mulparam) {
	// 最多同时允许100台配置的下载, 超过100台使其等待
	wgp := glimit.New(100)
	for secName, mulP := range sws {
		if err := createDirIfNotExist(secName); err != nil {
			logrus.Errorln(err)
			return
		}
		wgp.Add(len(mulP.profiles))
		for _, profile := range mulP.profiles {
			go func(sn string, pf param) {
				defer wgp.Done()
				if err := s.downloadFile(sn, pf); err != nil {
					logrus.Errorln(err)
				}
			}(secName, profile)
		}
	}
	wgp.Wait()
}

// downloadFile 从ftp下载文件
func (s *swbka) downloadFile(secName string, profile param) error {
	ip := profile.ip
	retErr := func(ip string, err error) error {
		return fmt.Errorf("address:%s reason:%s", ip, err.Error())
	}
	c, err := ftp.Dial(ip, ftp.DialWithTimeout(6*time.Second))
	if err != nil {
		return retErr(ip, err)
	}
	defer func() {
		if err := c.Quit(); err != nil {
			logrus.Errorf("ftp client quit failed: %v address: %s\n", err, ip)
		}
	}()
	// login ftp server
	err = c.Login(profile.username, profile.password)
	if err != nil {
		return retErr(ip, err)
	}
	// retrieve file content
	for _, fName := range profile.target {
		r, err := c.Retr(fName)
		if err != nil {
			continue
		}
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			return retErr(ip, err)
		}
		now := time.Now().Format(timeFormat)
		filename := joinString("/", secName, joinString("_", profile.deviceName, strings.Replace(profile.ip, ":", "_", -1), fName))
		err = s.saveFile(buf, filename, now)
		if err != nil {
			return retErr(ip, err)
		}
		if err := r.Close(); err != nil {
			if strings.Contains(err.Error(), "Entering Passive Mode") {
				logrus.Info("no such file: "+fName + " at " + profile.ip)
			}
		}
	}
	return nil
}

func (s *swbka) saveFile(data []byte,
	filename string, atTime string) error {
	if len(data) == 0 {
		return errors.New("data is empty")
	}
	// 无论是否有改动，都将推送至云盘
	if err := s.pushConfig2Webdav(
		data,
		joinPath("networkDeviceCFG",
			atTime, filename)); err != nil {
		return err
	}
	hash := checkSum(data)
	if hash == "" {
		return fmt.Errorf("failed to hash: %s\n", filename)
	}
	if fName, ok := s.sumMap.Load(hash); ok {
		return fmt.Errorf("cfg is not changed: %s\n", fName)
	}
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	// append filename md5 to cfg.sum
	file, err := os.OpenFile(s.sumFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logrus.Errorln(err.Error())
		}
	}()
	_, err = file.WriteString(fmt.Sprintf("%s %s\n", filename, hash))
	if err != nil {
		return err
	}
	return nil
}

func (s *swbka) readSumFile() error {
	// create cfg.sum file if it does not exist
	if _, err := os.Stat(s.sumFile); os.IsNotExist(err) {
		if _, err := os.Create(s.sumFile); err != nil {
			return err
		}
	}
	file, err := os.Open(s.sumFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logrus.Errorln(err)
		}
	}()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		if len(line) == 2 {
			s.sumMap.Store(line[1], line[0])
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	swb := &swbka{
		sumFile: "./cfg.sum",
		sumMap:  sync.Map{},
		wd:      gowebdav.NewClient(*webdavURI, *webdavUSER, *webdavPWD),
	}
	err := swb.readSumFile()
	if err != nil {
		logrus.Fatal(err)
	}
	ret, err := swb.readConfig(*configPath)
	if err != nil {
		logrus.Fatal(err)
	}
	swb.downloadSwitchCfg(ret)
}
