package adtools

import (
	"crypto/tls"

	"github.com/go-ldap/ldap/v3"
)

const (
	BaseDN = "dc=0x1un,dc=io"
)

type adConn struct {
	Conn *ldap.Conn
}

// buser => bind user, bpass => bind password
func NewADConn(url, buser, bpass string) (*adConn, error) {
	l, err := ldap.DialURL(url)
	if err != nil {
		return nil, err
	}
	if err := l.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
		return nil, err
	}
	if err := l.Bind(buser, bpass); err != nil {
		return nil, err
	}
	return &adConn{
		Conn: l,
	}, nil
}
