package adtools

import "fmt"

// query options
const (
	CnFilter                  = "(CN=%s)"
	CnFuzzyFilter             = "(CN=*%s*)"
	SAMAccountNameFilter      = "(sAMAccountName=%s)"
	SAMAccountNameFuzzyFilter = "(sAMAccountName=*%s*)"
	OuWithoutDefaultOUFilter  = "(&(objectClass=organizationalUnit)(ou=%s)(!(OU=Domain Controllers)))"
	AllUserFilter             = "(&(objectClass=User)(objectCategory=Person))"
)

var (
	DisabledFlag    = StringListWrap("514")
	EnabledFlag     = StringListWrap("512")
	ObjectClassBase = StringListWrap("top", "person", "organizationalPerson", "user")
)

func Ft(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
