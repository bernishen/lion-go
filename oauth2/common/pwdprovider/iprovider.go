package pwdprovider

import (
	"github.com/bernishen/exception"
)

// IPwdProvider :  the calculate method of the provider is defined.
type iPwdProvider interface {
	Compute(string) (string, *exception.Exception)
}
