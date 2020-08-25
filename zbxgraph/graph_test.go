package zbxgraph

import (
	"fmt"
	"testing"
)

func TestGethost(t *testing.T) {
	s := NewZbxGraph("http://10.100.100.150/api_jsonrpc.php", "Admin", "goodluck@123")
	res, err := s.ListGroup("滴滴")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}
