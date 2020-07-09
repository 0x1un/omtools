package adtools

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

var (
	UserIsAlreadyExsist = func(a interface{}) error { return fmt.Errorf("user is already exsists: %s", a) }
	IdenticalTarget     = func(a interface{}) error { return fmt.Errorf("have indentical target: %s", a) }
)

type Failed struct {
	Errors []error
}

// UserProfile is the necessary for create user
type UserInfo struct {
	Username           string   // not null
	Cn                 string   // not null
	Org                string   // not null
	Description        []string // not null
	UnicodePwd         []string // not null
	ObjectClass        []string // not null
	UserAccountControl []string // 514 is disabled, 512 is activate
	DisplayName        []string // not null
	SAMAccountName     []string // not null
}

func GenAttribute(profile UserInfo) *ldap.AddRequest {
	newreq := fmt.Sprintf("cn=%s,%s,"+BaseDN, profile.Cn, profile.Org)
	sqlInsert := ldap.NewAddRequest(newreq, nil)
	sqlInsert.Attribute("objectClass", profile.ObjectClass)
	sqlInsert.Attribute("cn", StringListWrap(profile.Cn))
	sqlInsert.Attribute("userAccountControl", profile.UserAccountControl)
	sqlInsert.Attribute("displayName", profile.DisplayName)
	sqlInsert.Attribute("unicodePwd", profile.UnicodePwd)
	sqlInsert.Attribute("sAMAccountName", profile.SAMAccountName)
	sqlInsert.Attribute("description", profile.Description)
	sqlInsert.Attribute("pwdLastSet", StringListWrap("0"))
	return sqlInsert
}
