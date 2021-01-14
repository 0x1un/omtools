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
	"github.com/0x1un/boxes/chatbot"
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
	dateFormat = "2006-01-02--15_04_05"
	configPath = flag.String("c", "/etc/swbka/profile.ini", "config file location")
)

type param struct {
	ip         string
	username   string
	password   string
	target     []string
	deviceName string
}

type general struct {
	dingNotifyAll bool
	smtpPort      int
	projectName   string
	pubUser       string
	pubPass       string
	pubPort       string
	webdavURL     string
	webdavUSER    string
	webdavPWD     string
	profilePATH   string
	smtpServer    string
	smtpUSER      string
	smtpPWD       string
	smtpFROM      string
	smtpTO        []string
	pubTarget     []string
	dingTokens    []string
	dingAtUsers   []string
}

type mulparam struct {
	profiles []param
}

// swbka: 入口结构
type swbka struct {
	sumFile    string
	sumMap     sync.Map
	wd         *gowebdav.Client
	defaultCFG general
	failed     []error
	total      int
	filesCount int
	davData    sync.Map
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
	func(switches map[string]mulparam) {
		for _, val := range switches {
			s.filesCount += len(val.profiles)
		}
	}(sws)
	// 最多同时允许100台配置的下载, 超过100台使其等待
	wgp := glimit.New(100)
	lock := sync.Mutex{}
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
					if er := err.Error(); !strings.Contains(er, "not changed") {
						lock.Lock()
						s.failed = append(s.failed, err)
						defer lock.Unlock()
					}
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
		return fmt.Errorf("address: %s reason: %s", ip, err.Error())
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
			logrus.Errorln(err)
			continue
		}
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			return retErr(ip, err)
		}
		filename := joinString("/", secName,
			joinString("-", profile.deviceName, strings.Replace(profile.ip,
				":",
				"-",
				-1),
				fName))
		// 存储配置数据
		s.davData.Store(filename, buf)
		err = s.saveFile(buf, filename)
		if err != nil {
			return retErr(ip, err)
		}
		if err := r.Close(); err != nil {
			if strings.Contains(err.Error(), "Entering Passive Mode") {
				logrus.Info("no such file: " + fName + " at " + profile.ip)
			}
		}
	}
	return nil
}

func (s *swbka) saveFile(data []byte, filename string) error {
	if len(data) == 0 {
		return errors.New("data is empty")
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

// Impl 具体调用的入口功能
func Impl() error {
	// 解析命令行参数
	flag.Parse()
	swb := &swbka{
		sumFile: "./cfg.sum",
		sumMap:  sync.Map{},
	}
	err := swb.readSumFile()
	if err != nil {
		return err
	}
	ret, err := swb.readConfig(*configPath)
	if err != nil {
		return err
	}
	// 初始化webdav客户端，并进行配置文件下载
	swb.wd = gowebdav.NewClient(swb.defaultCFG.webdavURL, swb.defaultCFG.webdavUSER, swb.defaultCFG.webdavPWD)
	swb.downloadSwitchCfg(ret)
	filesBuffer := files{}
	now := time.Now()
	nowTime := now.Format(dateFormat)
	// 将数据推送至webdav
	swb.davData.Range(func(filename, data interface{}) bool {
		filePath := strings.Split(filename.(string), "/")
		fileSection := func(a []string) string {
			if len(a) >= 2 {
				return a[0]
			}
			return ""
		}(filePath)
		cfgName := func(a []string) string {
			if len(a) >= 2 {
				return a[1]
			}
			return ""
		}(filePath)

		err := swb.pushConfig2Webdav(data.([]byte), joinString("/", "/networkDeviceCFG", swb.defaultCFG.projectName, fileSection, nowTime, cfgName))
		if err != nil {
			swb.failed = append(swb.failed, err)
		}
		filesBuffer = append(filesBuffer, file{
			filename: filename.(string),
			fileData: data.([]byte),
		})
		return true
	})
	downloadedCount := len(filesBuffer)
	zipDirectory := "./zip_packages/"
	if _, err := os.Stat(zipDirectory); os.IsNotExist(err) {
		if err := os.Mkdir(zipDirectory, 0644); err != nil {
			return err
		}
	}
	attachFile := zipDirectory + swb.defaultCFG.projectName + "_" + nowTime + ".zip"
	if err := filesBuffer.zip(attachFile); err != nil {
		return err
	}
	// 写入提示信息
	strBuffer := strings.Builder{}
	var botSendContent string
	if failedLen := len(swb.failed); failedLen != 0 {
		for _, er := range swb.failed {
			strBuffer.WriteString(er.Error() + "\n\n\n")
		}
		botSendContent = fmt.Sprintf("<font color=#003153>%s</font>\n\n<font color=#1E90FF>[%s]</font> ➤ 本次备份进度: <font color=#ff0000>%d/%d</font>\n\n但出现过错误，如下：\n\n<font color=#ff0000>%s</font>\n", now.Format(time.RFC3339), swb.defaultCFG.projectName, swb.total-failedLen, swb.total, strBuffer.String())
	} else {
		botSendContent = fmt.Sprintf("<font color=#003153>%s</font>\n\n<font color=#1E90FF>[%s]</font> ➤ 本次备份: <font color=#00ffff>%d/%d</font>\n\n", now.Format(time.RFC3339), swb.defaultCFG.projectName, swb.total, swb.total)
	}

	botSendContent = fmt.Sprintf("%s\n\n<font color=#ff0000>本次备份总共 %d 次, 有效 %d 次</font>\n", botSendContent, swb.filesCount, downloadedCount)

	// 进行钉钉告警
	chatbot.Send(
		swb.defaultCFG.dingTokens,
		swb.defaultCFG.dingAtUsers,
		swb.defaultCFG.dingNotifyAll,
		botSendContent, "backup message:")
	// 发送邮件存档
	sender := NewSMTPSender(swb.defaultCFG.smtpServer, swb.defaultCFG.smtpPort, swb.defaultCFG.smtpUSER, swb.defaultCFG.smtpPWD)
	msg := sender.writeMessage(botSendContent, swb.defaultCFG.smtpFROM, attachFile, swb.defaultCFG.projectName+" 配置存档", swb.defaultCFG.smtpTO...)
	if err := sender.SendToMail(msg); err != nil {
		return err
	}
	return nil
}
