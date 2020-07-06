package adtools

import (
	"fmt"
	"log"
	"testing"

	"github.com/go-ldap/ldap/v3"
)

const (
	LDAP     = "ldap://172.19.2.10"
	username = "administrator@0x1un.io"
	password = "gdlk@123"
)

var (
	conn *adConn
)

func init() {
	con, err := NewADConn(LDAP, username, password)
	if err != nil {
		log.Fatal(err)
	}
	conn = con
}

func TestPreReadFile(t *testing.T) {
	res := PreReadFile("./example2.csv")
	fmt.Println(res)
}

func TestAddUserMultiple(t *testing.T) {
	conn.AddUserMultiple("./example.csv", "ou=o1,ou=om", false)
}

func TestBind(t *testing.T) {
	conn.CheckAccount("aumujun", "jhrz@123..")
}

func TestResetPassword(t *testing.T) {
	if err := conn.ResetPasswd("aumujun", "goodluck@123", "ou=o1,ou=om"); err != nil {
		t.Fatal(err)
	}
}
func TestQueryUser(t *testing.T) {
	res, err := conn.QueryUser(OuWithoutDefaultOUFilter, false)
	if err != nil {
		t.Fatal(err)
	}
	res.PrettyPrint(2)
}

func TestDel(t *testing.T) {
	err := conn.DelUser("zhangqing", "ou=om")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddUser(t *testing.T) {
	err := conn.AddUser("zhangqing", "zhangjun", "ou=o1,ou=om", "goodluck@123", "null", false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestConn(t *testing.T) {
	searchRequest := ldap.NewSearchRequest(
		BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(CnFilter, "zhangqing"),
		StringListWrap(""), nil,
	)
	sr, err := conn.Conn.Search(searchRequest)
	if err != nil {
		t.Error(err)
	}
	sr.PrettyPrint(4)
}
