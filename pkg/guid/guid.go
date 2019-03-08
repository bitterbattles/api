package guid

import "github.com/rs/xid"

// New returns a new GUID
func New() string {
	return xid.New().String()
}
