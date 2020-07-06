package adtools

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/qiniu/iconv"
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/unicode"
)

// EncodePwd string to unicode string
func EncodePwd(pwd string) string {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM) // 使用小端编码
	pwdEncoded, err := utf16.NewEncoder().String("\"" + pwd + "\"")
	if err != nil {
		return ""
	}
	return pwdEncoded
}

// StringListWrap convert string to []string
func StringListWrap(s ...string) []string {
	ss := make([]string, 0)
	ss = append(ss, s...)
	return ss
}

func PreReadFile(path string) [][]string {
	recordsOrigin, err := readCSVFile(path)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	records := make([][]string, 0)
	for _, record := range recordsOrigin {
		records = append(records, removeDuplicateElem(record))
	}
	return records
}

func removeDuplicateElem(s []string) []string {
	chkMap := make(map[string]struct{})
	array := make([]string, 0)
	for _, item := range s {
		if _, ok := chkMap[item]; !ok {
			chkMap[item] = struct{}{}
			array = append(array, item)
		}
	}

	return array
}

func detectedEncoding(bytes []byte) (*chardet.Result, error) {
	detector := chardet.NewTextDetector()
	res, err := detector.DetectBest(bytes)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func guessEncodeType(file *os.File) (string, error) {
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	result, err := detectedEncoding(bytes)
	if err != nil {
		return "", err
	}
	if result.Confidence >= 10 {
		return result.Charset, nil
	}
	return "", fmt.Errorf("suspect charset: %s", result.Charset)
}

func convertCharset(to, from string, file *os.File) (*iconv.Reader, error) {
	cd, err := iconv.Open(to, from)
	if err != nil {
		return nil, err
	}
	defer cd.Close()
	r := iconv.NewReader(cd, file, 0)
	return r, nil
}

func readCSVFile(path string) (records [][]string, err error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	encodeType, err := guessEncodeType(f)
	if err != nil {
		return nil, fmt.Errorf("failed guess file charset: %s", err)
	}
	f.Seek(0, 0)

	if encodeType != "UTF-8" {
		cd, err := iconv.Open("utf-8", "gbk")
		if err != nil {
			return nil, err
		}
		defer cd.Close()
		r := iconv.NewReader(cd, f, 0)
		records, err = csv.NewReader(r).ReadAll()
		if err != nil {
			return nil, err
		}
		return records, nil
	}
	records, err = csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func Jdisable(flags bool) []string {
	if flags {
		return DisabledFlag
	}
	return EnabledFlag
}

// GetErrorCode matching the error number from string
// give a string: LDAP Result Code 68 "Entry Already Exists": 00000524: UpdErr: DSID-031A11E2, problem 6005 (ENTRY_EXISTS), data 0
// return 68
func GetErrorCode(msg string) uint16 {

	return 0
}
