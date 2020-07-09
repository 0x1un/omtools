package adtools

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/go-ldap/ldap/v3"
	"github.com/xlab/treeprint"
)

const (
	LDAP       = "ldap://172.19.2.10"
	username   = "administrator@0x1un.io"
	password   = "gdlk@123"
	baseDomain = "0x1un.io"
)

var (
	conn    *adConn
	already []string
	tree    = treeprint.New()
)

func init() {
	con, err := NewADConn(LDAP, username, password)
	if err != nil {
		log.Fatal(err)
	}
	conn = con
}

var (
	linked = NewLinked()
	flag   = false
)

func recur(searchBase string) {
	var query = conn.QueryUser
	if len(strings.TrimSpace(searchBase)) == 0 {
		fmt.Println("searchBase cannot be empty")
		return
	}
	// 从根向下搜索节点，深度为 1
	res, err := query(searchBase, Ft(OuWithoutDefaultOUFilter, "*"), ldap.ScopeSingleLevel)
	if err != nil {
		log.Fatal(err)
	}

	if !flag {
		for _, node := range res.Entries {
			tree.AddBranch(node.DN)
			linked.Push(element(node.DN))
		}
		flag = true
	}

	// 从搜索到的节点中遍历，将每个节点再次进行深度为1的搜索
	for _, node := range res.Entries {
		idx := linked.Search(element(strings.Join(strings.Split(node.DN, ",")[1:], ",")))
		if idx >= 0 {
			tre := tree.FindByValue(string(linked.Get(uint(idx)).Data))
			if tre != nil {
				tre.AddBranch(node.DN)
			}

		}
		recur(node.DN)
	}
}

func TestQueryDepth(t *testing.T) {
	tree = tree.AddBranch(BaseDN)
	recur(BaseDN)
	res, err := conn.QueryUser(BaseDN, Ft(OuWithoutDefaultOUFilter, "*"), ldap.ScopeWholeSubtree)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range res.Entries {
		tre := tree.FindByValue(strings.Split(v.DN, ",")[1:])
		if tre != nil {
			continue
		}
	}
	println(tree.String())

}

func TestLinkedList(t *testing.T) {
	link := NewLinked()
	for i := 0; i <= 26; i++ {
		link.Push(element(fmt.Sprintf("%c", 'a'+i)))
	}
	link.Print()
	println(link.Get(0).Data)
	println(link.Get(25))
	println(link.Search("a"))
	println(link.Search("b"))
	println(link.Search("c"))
	println(link.Search("d"))
	println(link.Search("z"))
	println(link.Len())
}

func TestTemp(t *testing.T) {
	failed := conn.MoveUserMultiple("testfiles/example2.csv", "o1")
	if len(failed.Errors) > 0 {
		PrintlnList(failed.Errors)
	}
}

func TestMoveUser(t *testing.T) {
	err := conn.MoveUser("zj", "cdjh-al")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDelUserMultiple(t *testing.T) {
	failed := conn.DelUserMultiple("testfiles/example2.csv", "ou=o1,ou=om")
	if fd := failed.Errors; len(fd) > 0 {
		PrintlnList(fd)
	}
}
func TestPreReadFile(t *testing.T) {
	res := PreReadFile("testfiles/example2.csv")
	fmt.Println(res)
}

func TestAddUserMultiple(t *testing.T) {
	failed := conn.AddUserMultiple("testfiles/example2.csv", "ou=o1,ou=om", false)
	if fd := failed.Errors; len(fd) > 0 {
		PrintlnList(fd)
	}
}

func TestBind(t *testing.T) {
	conn.CheckAccount("aumujun", "jhrz@123..")
}

func TestResetPassword(t *testing.T) {
	if err := conn.ResetPasswd("aumujun", "goodluck@123", "ou=o1,ou=om"); err != nil {
		t.Fatal(err)
	}
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
