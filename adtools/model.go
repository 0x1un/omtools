package adtools

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

// UserProfile is the necessary for create user
type UserInfo struct {
	Username           string
	Cn                 string
	Org                string
	Description        []string
	UnicodePwd         []string
	ObjectClass        []string
	UserAccountControl []string // 514 activate
	DisplayName        []string
	SAMAccountName     []string
}

func GenAttribute(profile UserInfo) *ldap.AddRequest {
	newreq := fmt.Sprintf("cn=%s,ou=%s,"+BaseDN, profile.Cn, profile.Org)
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
