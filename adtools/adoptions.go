package adtools

import "fmt"

// query options
const (
	CnFilter                  = "(CN=%s)"
	SAMAccountNameFilter      = "(sAMAccountName=%s)"
	SAMAccountNameFuzzyFilter = "(sAMAccountName=*%s*)"
	OuWithoutDefaultOUFilter  = "(&(objectClass=organizationalUnit)(!(OU=Domain Controllers)))"
)

var (
	DisabledFlag    = StringListWrap("514")
	EnabledFlag     = StringListWrap("512")
	ObjectClassBase = StringListWrap("top", "person", "organizationalPerson", "user")
)

func Ft(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
