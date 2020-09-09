package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var (
	timeFormat = "20060102"
)

type param struct {
	ip       string
	username string
	password string
	target   []string
}

type mulparam struct {
	profiles []param
}

// swbka: 入口结构
type swbka struct {
	sumFile string
	sumMap  map[string]string
}

// readConfig 读取配置文件
func (*swbka) readConfig(path string) (map[string]mulparam, error) {
	mp := make(map[string]mulparam)
	cfg, err := ini.LoadSources(ini.LoadOptions{SkipUnrecognizableLines: true, IgnoreInlineComment: true}, path)
	if err != nil {
		return nil, err
	}
	pubUser := cfg.Section("general").Key("pub_user").String()
	pubPass := cfg.Section("general").Key("pub_pass").String()
	pubTarget := cfg.Section("general").Key("pub_target").String()
	if pubTarget == "" {
		pubTarget = "startup.cfg"
	}
	for _, v := range cfg.Sections() {
		name := v.Name()
		if name == "general" || name == "DEFAULT" {
			continue
		}
		m := mulparam{}
		for ip, loginStr := range v.KeysHash() {
			strList := strings.Split(loginStr, ",")
			if len(strList) == 3 {
				m.profiles = append(m.profiles, param{
					ip:       ip,
					username: strList[0],
					password: strList[1],
					target:   []string{strList[2]},
				})
			} else {
				m.profiles = append(m.profiles, param{
					ip:       ip,
					username: pubUser,
					password: pubPass,
					target:   strings.Split(pubTarget, ","),
				})
			}
		}
		mp[name] = m
	}
	return mp, nil
}

// Deprecated: use downloadSwitchCfg instead.
func downloadFunc(s *swbka, wg *sync.WaitGroup, secName string, mulP mulparam) {
	if _, err := os.Stat(secName); os.IsNotExist(err) {
		err := os.Mkdir(secName, 0644)
		if err != nil {
			logrus.Fatal(err)
		}
	}
	for _, profile := range mulP.profiles {
		go func(sn string, pf param) {
			// err := s.downloadFile(sn, pf)
			err := s.downloadFileMock(sn, pf)
			wg.Done()
			if err != nil {
				logrus.Errorln(err)
			}
		}(secName, profile)
	}
}

func (s *swbka) downloadFileMock(sn string, p param) error {
	fmt.Println(p.ip)
	time.Sleep(2 * time.Second)
	return nil
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0644); err != nil {
			return err
		}
	}
	return nil
}

// downloadSwitchCfg 批量下载交换机配置文件
func (s *swbka) downloadSwitchCfg(sws map[string]mulparam) {
	wgp := sync.WaitGroup{}
	for secName, mulP := range sws {
		if err := createDirIfNotExist(secName); err != nil {
			logrus.Fatal(err)
		}

		wgp.Add(len(mulP.profiles))

		for _, profile := range mulP.profiles {
			go func(sn string, pf param) {
				if err := s.downloadFile(sn, pf); err != nil {
					logrus.Errorln(err)
				}
				wgp.Done()
			}(secName, profile)
		}
	}
	wgp.Wait()
}

// downloadFile 从ftp下载文件
func (s *swbka) downloadFile(secName string, profile param) error {
	ip := profile.ip + ":21"
	retErr := func(ip string, err error) error {
		return fmt.Errorf("%s ➨ %s", ip, err)
	}
	c, err := ftp.Dial(ip, ftp.DialWithTimeout(6*time.Second))
	if err != nil {
		return retErr(ip, err)
	}
	// login ftp server
	err = c.Login(profile.username, profile.password)
	if err != nil {
		return retErr(ip, err)
	}
	// retrieve file content
	for _, fname := range profile.target {
		r, err := c.Retr(fname)
		if err != nil {
			continue
		}
		defer r.Close()

		buf, err := ioutil.ReadAll(r)
		if err != nil {
			return retErr(ip, err)
		}
		now := time.Now().Format(timeFormat)
		filename := secName + "/" + profile.ip + "_" + now + "_" + fname
		err = s.saveFile(buf, filename)
		if err != nil {
			return retErr(ip, err)
		}
	}
	if err := c.Quit(); err != nil {
		return retErr(ip, err)
	}
	return nil
}

func (s *swbka) saveFile(data []byte, filename string) error {
	hash := checkSum(data)
	if hash == "" {
		return fmt.Errorf("failed to hash: %s\n", filename)
	}
	if fname, ok := s.sumMap[hash]; ok {
		return fmt.Errorf("cfg is not changed: %s\n", fname)
	}
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return err
	}
	// append filename md5 to cfg.sum
	file, err := os.OpenFile(s.sumFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
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
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		if len(line) == 2 {
			s.sumMap[line[1]] = line[0]
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// checkSum 获取一段文本内容的hash值
func checkSum(content []byte) string {
	sha := md5.New()
	_, err := sha.Write(content)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(sha.Sum(nil))
}

func main() {
	start := time.Now()
	swb := &swbka{
		sumFile: "./cfg.sum",
		sumMap:  make(map[string]string),
	}
	err := swb.readSumFile()
	if err != nil {
		logrus.Fatal(err)
	}
	ret, err := swb.readConfig("./profile.ini")
	if err != nil {
		logrus.Fatal(err)
	}
	swb.downloadSwitchCfg(ret)
	fmt.Printf("time count: %ds\n", time.Now().Second()-start.Second())
}
