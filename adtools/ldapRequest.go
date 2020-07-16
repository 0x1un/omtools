// this file included the ldap request methods

package adtools

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/manifoldco/promptui"
)

func NewADTools(url, buser, bpass string) (ADTooller, error) {
	conn, err := NewADConn(url, buser, bpass)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type ADTooller interface {
	// 添加单个用户
	AddUser(disName, username, orgName, loginPwd, description string, disabled bool) error
	// 删除单个用户
	DelUser(disName, ouPath string) error
	// 从主dn中查询
	QueryUserFromBaseDN(filter string) (*ldap.SearchResult, error)
	// 从指定dn中查询
	QueryUser(dn, filter string, scope int) (*ldap.SearchResult, error)
	// 重置密码
	ResetPasswd(uname, passwd, ouPath string) error
	// 检测账户状态，是否可用
	CheckAccount(username, password string)
	// 从csv文件中添加用户
	AddUserMultiple(importPath, orgName string, disabled bool) Failed
	// 从csv文件中删除用户
	DelUserMultiple(path, orgName string) Failed
	// 移动用户
	MoveUser(from, to string) error
	// 移动csv文件中的用户
	MoveUserMultiple(path, to string) Failed
	// 内置可暴露ldap链接
	BuiltinConn() *ldap.Conn
	// 移动用户至绝对路径中
	MoveUserAbsPath(from, to string) error
	// GetUserInfoFormat(user string) (string, error)
	// 获取用户信息并生成位可视table
	GetUserInfoTable(user string) (string, error)
	// 导出csv模板
	ExportTemplate(target string)
	// 禁用或启用用户
	ChangeUserStatus(dn string, disabled bool) error
	// 解锁或锁定用户
	UnlockUser(dn string) error
}

func (c *adConn) BuiltinConn() *ldap.Conn {
	return c.Conn
}

// ChangeUserStatus 需要完整的dn路径，以及更改的状态
func (c *adConn) ChangeUserStatus(dn string, disabled bool) error {
	modifyReq := ldap.NewModifyRequest(dn, nil)
	modifyReq.Replace("userAccountControl", Jdisable(disabled))
	return c.Conn.Modify(modifyReq)
}

func (c *adConn) UnlockUser(dn string /*, locked bool*/) error {
	modifyReq := ldap.NewModifyRequest(dn, nil)
	modifyReq.Replace("lockoutTime", StringListWrap("0"))
	return c.Conn.Modify(modifyReq)

}
func (c *adConn) ExportTemplate(target string) {
	data := []byte(`姓名,批次,域账号,密码
吕布,1001,lvbu,abc@12345
刘备,1001,liubei,abc@12345`)
	err := ioutil.WriteFile(target, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
func (c *adConn) RemoveExpireComputerOnTime() {}

// AddUser add a single user to specify ou
// orgName mut be "ou=01,ou=om"
/*
	0x1un.io
		om
			o1
*/
func (c *adConn) AddUser(disName, username, orgName, loginPwd, description string, disabled bool) error {
	attr := GenAttribute(UserInfo{
		Username:           username,
		Cn:                 disName,
		Org:                orgName,
		SAMAccountName:     StringListWrap(username),
		UserAccountControl: Jdisable(disabled),
		ObjectClass:        ObjectClassBase,
		UnicodePwd:         StringListWrap(EncodePwd(loginPwd)),
		Description:        StringListWrap(description),
		DisplayName:        StringListWrap(disName),
	})
	err := c.Conn.Add(attr)
	if err != nil {
		if ldap.IsErrorWithCode(err, 68) {
			return fmt.Errorf("%s is already exists\n", disName)
		}
		ldapErr, ok := err.(*ldap.Error)
		if !ok {
			return fmt.Errorf("unkown error: %s", err.Error())
		}
		return fmt.Errorf("dn: %s, msg: %s\n", disName, ldapErr.Error())
	}
	fmt.Printf("added user: %s\n", disName)
	return nil
}

// DelUser given an existing account to delete
// ouPath must be a detailed path, ignore the domain name
// for example: "ou=o1,ou=om"
/*
	0x1un.io
		om
			o1
*/
func (c *adConn) DelUser(disName, ouPath string) error {
	delDn := fmt.Sprintf("CN=%s,%s,%s", disName, ouPath, BaseDN)
	delReq := ldap.NewDelRequest(delDn, nil)
	err := c.Conn.Del(delReq)
	if err != nil {
		ldapErr, ok := err.(*ldap.Error)
		if !ok {
			return fmt.Errorf("unkown error: %s", err.Error())
		}
		if ldap.IsErrorWithCode(err, 32) {
			return fmt.Errorf("del dn: %s, %s", disName, ldap.LDAPResultCodeMap[ldapErr.ResultCode])
		}
		return fmt.Errorf("del dn: %s, %s", disName, ldap.LDAPResultCodeMap[ldapErr.ResultCode])
	}
	return nil
}

// QueryUser get single user information
func (c *adConn) QueryUserFromBaseDN(filter string) (*ldap.SearchResult, error) {
	return c.QueryUser(BaseDN, filter, ldap.ScopeWholeSubtree)
}

func (c *adConn) QueryUser(dn, filter string, scope int) (*ldap.SearchResult, error) {
	searchReq := ldap.NewSearchRequest(
		dn,
		scope, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		StringListWrap("*"), nil)
	res, err := c.Conn.Search(searchReq)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ResetPasswd to reset exsisting account password
func (c *adConn) ResetPasswd(uname, passwd, ouPath string) error {
	modifyReq := ldap.NewModifyRequest(Ft("CN=%s,%s,"+BaseDN, uname, ouPath), nil)
	modifyReq.Replace("unicodePwd", StringListWrap(EncodePwd(passwd)))
	return c.Conn.Modify(modifyReq)
}

// CheckAccount to checking account is it available
func (c *adConn) CheckAccount(username, password string) {
	err := c.Conn.Bind(username, password)

	if err != nil {
		fmt.Printf("failed bind to user: %s", err)
		return
	}
	fmt.Println("bind successfully")
}

// AddUserMultiple import user from csv file
// the format must be fllows:
// 姓名,批次,域账号,密码
// 张三,1001,zhangsan,jnk@123.
// 只返回失败的item
func (c *adConn) AddUserMultiple(importPath, orgName string, disabled bool) Failed {
	records := PreReadFile(importPath)
	failed := Failed{}
	// convert record to *ldap.AddRequest list
LOOP_1:
	// ignored head line
	for _, record := range records[1:] {
		// continue if record is empty
		for _, r := range record {
			if r == "" {
				continue LOOP_1
			}
		}
		err := c.AddUser(record[0], record[2], orgName, record[3], record[1], disabled)
		if err != nil {
			failed.Errors = append(failed.Errors, err)
		}
	}

	return failed
}

// DelUserMultiple import user from csv file and remove them
// the format must be fllows:
// 姓名,批次,域账号,密码
// 张三,1001,zhangsan,jnk@123.
// 只返回失败的item
func (c *adConn) DelUserMultiple(path, orgName string) Failed {
	failed := Failed{}
	records := PreReadFile(path)
	for _, record := range records[1:] {
		err := c.DelUser(record[0], orgName)
		if err != nil {
			failed.Errors = append(failed.Errors, err)
		}
	}
	return failed
}

// MoveUser move only
// for example: move Harry to o2 organization
// MoveUser("cn=Harry,ou=o1,ou=om,dc=0x1un,dc=io","ou=o2,ou=om,dc=0x1un,dc=io")
func (c *adConn) MoveUserAbsPath(from, to string) error {
	getFirst := func(s string) string {
		list := strings.Split(s, ",")
		if len(list) == 0 {
			return ""
		}
		return list[0]
	}
	fmt.Println(getFirst(from))
	modifyReq := ldap.NewModifyDNRequest(from, getFirst(from), true, to)
	err := c.Conn.ModifyDN(modifyReq)
	if err != nil {
		return err
	}
	return nil
}

// MoveUser move only, from must be a CN
// for example:
// MoveUser("Harry", "o2")
// 0x1un.io
// 	om
// 		o1
// 		o2
func (c *adConn) MoveUser(from, to string) error {
	fromFilter := Ft(CnFilter, from)
	toFilter := Ft(OuWithoutDefaultOUFilter, to)
	fromQuery, err := c.QueryUserFromBaseDN(fromFilter)
	if err != nil {
		return err
	}
	toQuery, err := c.QueryUserFromBaseDN(toFilter)
	if err != nil {
		return err
	}
	if len(fromQuery.Entries) == 0 {
		return NotFound(from)
	}
	if len(toQuery.Entries) == 0 {
		return NotFound(to)
	}
	// 目标有多个, 提醒用户去选择
	selectPath := toQuery.Entries[0].DN
	if entry := toQuery.Entries; len(entry) > 1 {
		dnn := []string{}
		for _, v := range entry {
			dnn = append(dnn, v.DN)
		}
		pSelect := promptui.Select{
			Label: "Select path",
			Items: dnn,
		}
		_, res, err := pSelect.Run()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		selectPath = res
	}
	// 判断user是否已在此组
	if list := strings.Split(fromQuery.Entries[0].DN, ","); len(list) > 0 {
		if strings.Join(list[1:], ",") == toQuery.Entries[0].DN {
			return UserIsAlreadyExsist(from)
		}
	}
	err = c.MoveUserAbsPath(fromQuery.Entries[0].DN, selectPath)
	if err != nil {
		return err
	}
	return nil
}

// MoveUserMultiple move user from csv file
// 姓名,批次,域账号,密码
// 张三,1001,zhangsan,jnk@123.
func (c *adConn) MoveUserMultiple(path, to string) Failed {
	records := PreReadFile(path)
	failed := Failed{}
	for _, v := range records[1:] {
		if err := c.MoveUser(v[0], to); err != nil {
			failed.Errors = append(failed.Errors, err)
		}
	}
	return failed
}
