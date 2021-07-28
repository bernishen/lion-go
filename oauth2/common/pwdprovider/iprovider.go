package pwdprovider

import (
	"github.com/Berni-Shen/lion-go/utils/exception"
)

// IPwdProvider :  the calculate method of the provider is defined.
type iPwdProvider interface {
	Compute(string) (string, *exception.Exception)
}
