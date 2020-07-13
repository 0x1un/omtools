package adtools

import (
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

// 查询一个用户的基本信息 并排列格式到一个buffer中
// func (c *adConn) GetUserInfoFormat(user string) (string, error) {
// 	res, err := c.QueryUserFromBaseDN(fmt.Sprintf("(|(CN=%s)(sAMAccountName=%s))", user, user))
// 	if err != nil {
// 		return "", err
// 	}
// 	buffer := strings.Builder{}
// 	buffer.WriteString(fmt.Sprintf("%s 基本信息\n", user))
// 	for _, v := range res.Entries {
// 		buffer.WriteString(Ft("\t姓名: %s\n", v.GetAttributeValue("cn")))
// 		buffer.WriteString(Ft("\t账户: %s\n", v.GetAttributeValue("sAMAccountName")))
// 		buffer.WriteString(Ft("\t创建时间: %s\n", v.GetAttributeValue("whenCreated")))
// 		buffer.WriteString(Ft("\t最后更改时间: %s\n", v.GetAttributeValue("whenChanged")))
// 		buffer.WriteString(Ft("\t用户激活状态: %s\n", v.GetAttributeValue("userAccountControl")))
// 		buffer.WriteString(Ft("\t密码错误次数: %s\n", v.GetAttributeValue("badPwdCount")))
// 		buffer.WriteString(Ft("\t最后登入时间: %s\n", v.GetAttributeValue("lastLogon")))
// 		buffer.WriteString(Ft("\t最后登出时间: %s\n", v.GetAttributeValue("lastLogoff")))
// 		buffer.WriteString(Ft("\t登入次数: %s\n", v.GetAttributeValue("logonCount")))
// 		buffer.WriteString(Ft("\t用户DN路径: %s\n", v.GetAttributeValue("distinguishedName")))
// 		buffer.WriteString(Ft("\t是否在下次登入时修改密码: %s\n", v.GetAttributeValue("pwdLastSet")))
// 	}
// 	return buffer.String(), nil
// }

func (c *adConn) GetUserInfoTable(user string) (string, error) {
	res, err := c.QueryUserFromBaseDN(Ft(UserFilter, user, user))
	if err != nil {
		return "", err
	}
	if len(res.Entries) == 0 {
		return "no record found", nil
	}
	tb := table.NewWriter()
	tb.AppendHeader(table.Row{
		"姓名", "账户", "创建时间",
		"最后更改时间", "用户激活状态",
		"密码错误次数", "最后登入时间",
		/*"最后登出时间",*/ "登入次数",
		"用户DN路径", "下次登入改密码"})
	for _, v := range res.Entries {
		tb.AppendRow([]interface{}{
			v.GetAttributeValue("cn"),
			v.GetAttributeValue("sAMAccountName"),
			parseDatetime2Humman(v.GetAttributeValue("whenCreated")),
			parseDatetime2Humman(v.GetAttributeValue("whenChanged")),
			func(a string) string {
				switch a {
				case "512":
					return "激活"
				case "514":
					return "禁用"
				case "66082":
					// Disabled, Password Doesn’t Expire & Not Required
					return "禁用"
				}
				return "unkown"
			}(v.GetAttributeValue("userAccountControl")),
			v.GetAttributeValue("badPwdCount"),
			func(a string) string {
				if strings.HasPrefix(a, "1601-01-01") {
					return "从未"
				}
				return a
			}(convertWinNTTime2Unix(v.GetAttributeValue("lastLogon"))),
			// v.GetAttributeValue("lastLogoff"),
			v.GetAttributeValue("logonCount"),
			v.GetAttributeValue("distinguishedName"),
			func(a string) string {
				if a == "0" {
					return "Yes"
				}
				return "No"
			}(v.GetAttributeValue("pwdLastSet")),
		})
		tb.AppendSeparator()
	}
	tb.SortBy([]table.SortBy{{Name: "姓名", Mode: table.AscNumeric}})
	return tb.Render(), nil
}

func convertWinNTTime2Unix(tm string) string {
	unum, err := strconv.ParseInt(tm, 10, 64)
	if err != nil {
		return "1970-01-01 00:00:00"
	}
	unum = (unum / 10000000) - 11644473600
	tim := time.Unix(unum, 0)
	return (tim.Format("2006-01-02 15:04:05"))
}

func parseDatetime2Humman(tm string) string {
	t, err := time.Parse("20060102150405.0Z", tm)
	if err != nil {
		return ""
	}
	return (t.Format("2006-01-02 15:04:05"))
}
