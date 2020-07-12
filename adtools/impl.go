package adtools

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// 查询一个用户的基本信息 并排列格式到一个buffer中
func (c *adConn) GetUserInfoFormat(user string) (string, error) {
	res, err := c.QueryUserFromBaseDN(fmt.Sprintf("(|(CN=%s)(sAMAccountName=%s))", user, user))
	if err != nil {
		return nil, err
	}
	buffer := strings.Builder{}
	buffer.WriteString(fmt.Sprintf("%s 基本信息\n", user))
	for _, v := range res.Entries {
		buffer.WriteString(Ft("\t姓名: %s\n", v.GetAttributeValue("cn")))
		buffer.WriteString(Ft("\t账户: %s\n", v.GetAttributeValue("sAMAccountName")))
		buffer.WriteString(Ft("\t创建时间: %s\n", v.GetAttributeValue("whenCreated")))
		buffer.WriteString(Ft("\t最后更改时间: %s\n", v.GetAttributeValue("whenChanged")))
		buffer.WriteString(Ft("\t用户激活状态: %s\n", v.GetAttributeValue("userAccountControl")))
		buffer.WriteString(Ft("\t密码错误次数: %s\n", v.GetAttributeValue("badPwdCount")))
		buffer.WriteString(Ft("\t最后登入时间: %s\n", v.GetAttributeValue("lastLogon")))
		buffer.WriteString(Ft("\t最后登出时间: %s\n", v.GetAttributeValue("lastLogoff")))
		buffer.WriteString(Ft("\t登入次数: %s\n", v.GetAttributeValue("logonCount")))
		buffer.WriteString(Ft("\t用户DN路径: %s\n", v.GetAttributeValue("distinguishedName")))
		buffer.WriteString(Ft("\t是否在下次登入时修改密码: %s\n", v.GetAttributeValue("pwdLastSet")))
	}

	return buffer.String(), nil
}

func (c *adConn) GetUserInfoTable(user string) (string, error) {
	res, err := c.QueryUserFromBaseDN(AllUserFilter)
	if err != nil {
		return "", err
	}
	tb := table.NewWriter()
	tb.AppendHeader(table.Row{
		"姓名", "账户", "创建时间",
		"最后更改时间", "用户激活状态",
		"密码错误次数", "最后登入时间",
		"最后登出时间", "登入次数",
		"用户DN路径", "下次登入改密码"})
	// 18位windows nt时间戳转换为unix时间戳
	// ( nt timestamp / 10000000 ) - 11644473600
	// (132389476541531163 / 10000000) -  11644473600
	for _, v := range res.Entries {
		tb.AppendRow([]interface{}{
			v.GetAttributeValue("cn"),
			v.GetAttributeValue("sAMAccountName"),
			v.GetAttributeValue("whenCreated"),
			v.GetAttributeValue("whenChanged"),
			v.GetAttributeValue("userAccountControl"),
			v.GetAttributeValue("badPwdCount"),
			v.GetAttributeValue("lastLogon"),
			v.GetAttributeValue("lastLogoff"),
			v.GetAttributeValue("logonCount"),
			v.GetAttributeValue("distinguishedName"),
			v.GetAttributeValue("pwdLastSet")})
		tb.AppendSeparator()
	}

	return tb.Render(), nil
}
