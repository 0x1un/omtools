package adtools

import "fmt"

// query options
const (
	CnFilter                  = "(CN=%s)"
	CnFuzzyFilter             = "(CN=*%s*)"
	SAMAccountNameFilter      = "(sAMAccountName=%s)"
	SAMAccountNameFuzzyFilter = "(sAMAccountName=*%s*)"
	OuWithoutDefaultOUFilter  = "(&(objectClass=organizationalUnit)(ou=%s)(!(OU=Domain Controllers)))"
	UserFilter                = "(&(objectClass=User)(objectCategory=Person)(|(CN=%s)(sAMAccountName=%s)))"
	LockedAllUserFilter       = "(&(objectCategory=Person)(objectClass=User)(lockoutTime>=1))"
)

var (
	DisabledFlag    = StringListWrap("514")
	EnabledFlag     = StringListWrap("512")
	ObjectClassBase = StringListWrap("top", "person", "organizationalPerson", "user")
)

func Ft(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
