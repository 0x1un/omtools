package main

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"strings"
)

func removeSpace(s []string) []string {
	re := []string{}
	for _, v := range s {
		if strings.TrimSpace(v) != "" {
			re = append(re, v)
		}
	}
	return re
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0644); err != nil {
			return err
		}
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

func joinPath(args ...string) string {
	return strings.Join(args, "/")
}

func joinString(sep string, args ...string) string {
	return strings.Join(args, sep)
}
