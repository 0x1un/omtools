// this file included the ldap request methods

package adtools

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
)

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
	if err := c.Conn.Add(attr); err != nil {
		return err
	}
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
	return c.Conn.Del(delReq)
}

// QueryUser get single user information
func (c *adConn) QueryUser(filter string, fuzzy bool) (*ldap.SearchResult, error) {
	searchReq := ldap.NewSearchRequest(
		BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		StringListWrap(""), nil)
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
func (c *adConn) AddUserMultiple(importPath, orgName string, disabled bool) error {
	records := PreReadFile(importPath)
	attributes := make([]*ldap.AddRequest, 0)
	// convert record to *ldap.AddRequest list
LOOP_1:
	for _, record := range records[1:] {
		for _, r := range record {
			if r == "" {
				continue LOOP_1
			}
		}
		attributes = append(attributes, GenAttribute(UserInfo{
			Username:           record[2],
			Cn:                 record[0],
			Org:                orgName,
			SAMAccountName:     StringListWrap(record[2]),
			UserAccountControl: Jdisable(disabled),
			ObjectClass:        StringListWrap("top", "person", "organizationalPerson", "user"),
			UnicodePwd:         StringListWrap(EncodePwd(record[3])),
			Description:        StringListWrap(record[1]),
			DisplayName:        StringListWrap(record[0]),
		}))
	}

	faileds := make([]string, 0)
	for _, reqAttr := range attributes {
		err := c.Conn.Add(reqAttr)
		if err != nil {
			if ldap.IsErrorWithCode(err, 68) {
				log.Printf("%s is already exists\n", reqAttr.DN)
				continue
			}
			ldapErr, ok := err.(*ldap.Error)
			if !ok {
				log.Printf("unkown error: %s\n", err)
				faileds = append(faileds, reqAttr.DN)
				continue
			}
			msg := ldap.LDAPResultCodeMap[ldapErr.ResultCode]
			log.Printf("dn: %s, msg: %s\n", reqAttr.DN, msg)
		}
	}
	if len(faileds) != 0 {
		fmt.Println("failed add dn:")
		for _, failed := range faileds {
			log.Println(failed)
		}
	}
	return nil
}
