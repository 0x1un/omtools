// this file included the ldap request methods

package adtools

import (
	"fmt"
	"log"
	"strings"

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
func (c *adConn) QueryUser(filter string) (*ldap.SearchResult, error) {
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
// 只返回失败的item
func (c *adConn) AddUserMultiple(importPath, orgName string, disabled bool) Failed {
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

	faileds := Failed{}
	for _, reqAttr := range attributes {
		err := c.Conn.Add(reqAttr)
		if err != nil {
			if ldap.IsErrorWithCode(err, 68) {
				log.Printf("%s is already exists\n", reqAttr.DN)
				continue
			}
			ldapErr, ok := err.(*ldap.Error)
			if !ok {
				faileds.Errors = append(faileds.Errors, fmt.Errorf(err.Error()+": %s", reqAttr.DN))
				continue
			}
			msg := ldap.LDAPResultCodeMap[ldapErr.ResultCode]
			log.Printf("dn: %s, msg: %s\n", reqAttr.DN, msg)
		}
	}
	return faileds
}

// DelUserMultiple import user from csv file and remove them
// the format must be fllows:
// 姓名,批次,域账号,密码
// 张三,1001,zhangsan,jnk@123.
// 只返回失败的item
func (c *adConn) DelUserMultiple(path, orgName string) Failed {
	faileds := Failed{}
	records := PreReadFile(path)
	for _, record := range records[1:] {
		err := c.DelUser(record[0], orgName)
		if err != nil {
			_, ok := err.(*ldap.Error)
			if !ok {
				faileds.Errors = append(faileds.Errors, err)
			} else if ldap.IsErrorWithCode(err, 32) {
				faileds.Errors = append(faileds.Errors, fmt.Errorf("no such object: %s", record[0]))
			}
			continue
		}
		fmt.Printf("del user: %s\n", record[0])
	}
	return faileds
}

// MoveUser move only
// for example: move Harry to o2 organization
// MoveUser("cn=Harry,ou=o1,ou=om,dc=0x1un,dc=io","ou=o2,ou=om,dc=0x1un,dc=io")
func (c *adConn) moveUserAbsPath(from, to string) error {
	getFirst := func(s string) string {
		list := strings.Split(s, ",")
		if len(list) == 0 {
			return ""
		}
		return list[0]
	}
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
	fromQuery, err := c.QueryUser(fromFilter)
	if err != nil {
		return err
	}
	toQuery, err := c.QueryUser(toFilter)
	if err != nil {
		return err
	}
	if len(fromQuery.Entries) == 0 || len(toQuery.Entries) == 0 {
		return fmt.Errorf("failed to query dn path: %s or %s", from, to)
	}
	// 判断user是否已在此组
	if list := strings.Split(fromQuery.Entries[0].DN, ","); len(list) > 0 {
		if strings.Join(list[1:], ",") == toQuery.Entries[0].DN {
			return UserIsAlreadyExsist(from)
		}
	}
	err = c.moveUserAbsPath(fromQuery.Entries[0].DN, toQuery.Entries[0].DN)
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
