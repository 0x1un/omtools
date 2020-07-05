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
	conn *ldap.Conn
)

func init() {
	con, err := NewADConn(LDAP, username, password)
	if err != nil {
		log.Fatal(err)
	}
	conn = con.Conn
}

func TestAddUser(t *testing.T) {
	user := GenAttribute(UserInfo{
		Username:           "akanami",
		Cn:                 "zhangjun",
		Org:                "om",
		SAMAccountName:     []string{"akanami"},
		UserAccountControl: []string{"512"},
		ObjectClass:        []string{"top", "person", "organizationalPerson", "user"},
		UnicodePwd:         []string{EncodePwd("goodluck@123")},
		Description:        []string{"1"},
		DisplayName:        []string{"zhangjun"},
	})
	err := conn.Add(user)
	if err != nil {
		t.Fatal(err)
	}
}

func TestConn(t *testing.T) {
	searchRequest := ldap.NewSearchRequest(
		BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(ObjectClassOrgbyUID, "administrator"),
		StringListWrap("dn"), nil,
	)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sr)

}
