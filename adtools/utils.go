package adtools

import "golang.org/x/text/encoding/unicode"

func EncodePwd(pwd string) string {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM) // 使用小端编码
	pwdEncoded, err := utf16.NewEncoder().String("\"" + pwd + "\"")
	if err != nil {
		return ""
	}
	return pwdEncoded
}

func StringListWrap(s ...string) []string {
	ss := make([]string, len(s))
	ss = append(ss, s...)
	return ss
}
