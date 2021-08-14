package pwdservice

import (
	"github.com/bernishen/lion-go/oauth2/common/pwdprovider"
	"github.com/bernishen/lion-go/oauth2/service/securityservice"
	"github.com/bernishen/exception"
)

// ConvertPwd : Converted the password format from page to database.
func ConvertPwd(clientID string, pwd string, pwdVersion int) (string, *exception.Exception) {
	pwdPlaintext, ex := securityservice.DecryptData(clientID, pwd)
	if ex != nil {
		return "", ex
	}

	ret, ex := pwdprovider.Compute(pwdVersion, pwdPlaintext)
	if ex != nil {
		return "", ex
	}

	return ret, nil
}
