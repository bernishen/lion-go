package pwdprovider

import (
	"github.com/bernishen/exception"
)

var pool [2]iPwdProvider

func init() {
	pool[0] = new(providerV1)
	pool[1] = new(providerV2)
}

// Compute  the text to get a cryphertext.
// 'pwdVersion' is the basis for provider.
//
// Exceptions:
//
// 1001: 'pwdVersion' can't less the 1.
//
// 1002: Found not the pwd provider.
//
// 2xxx: Has an exception in the provider computed.
func Compute(pwdVsersion int, text string) (string, *exception.Exception) {
	if pwdVsersion < 1 {
		return "", exception.NewException(exception.Error, 1001, "'pwdVersion' can't less the 1.")
	}
	if pwdVsersion > len(pool) {
		return "", exception.NewException(exception.Error, 1002, "'pwdVersion' is more than max value, the max is ")
	}
	p := pool[pwdVsersion-1]

	return p.Compute(text)
}

// NowPasswordVersion : Gets the current and latest password version.
func NowPasswordVersion() int {
	return len(pool)
}
