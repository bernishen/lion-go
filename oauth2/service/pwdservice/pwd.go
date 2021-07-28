package pwdservice

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/Berni-Shen/lion-go/utils/exception"
	"strings"

	"github.com/Berni-Shen/lion-go/oauth2/common/pwdprovider"
	"github.com/Berni-Shen/lion-go/oauth2/service/rsaservice"
)

// ConvertPwd : Converted the password format from page to database.
func ConvertPwd(clientID string, pwd string, pwdVersion int) (string, *exception.Exception) {
	rsaKey, _ := rsaservice.FindKeyByClient(clientID)

	var buffer bytes.Buffer
	buffer.WriteString(rsaKey)
	block, _ := pem.Decode(buffer.Bytes())
	if block == nil {
		return "", exception.NewException(exception.Error, 1001, "An exception occurred reading the 'PEM' key, the block can't null.")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", exception.NewException(exception.Error, 1002, "Parse the key block exception.["+err.Error()+"]")
	}

	buffer.Reset()
	buffer.WriteString(pwd)
	pwdNew, err := rsa.DecryptPKCS1v15(rand.Reader, key, buffer.Bytes())
	if err != nil {
		return "", exception.NewException(exception.Error, 1003, "An exception occurred encrypting.["+err.Error()+"]")
	}

	var builder strings.Builder
	for i := 0; i < len(pwdNew); i++ {
		builder.WriteByte(pwdNew[i])
	}

	ret, ex := pwdprovider.Compute(pwdVersion, builder.String())
	if ex != nil {
		return "", ex
	}

	return ret, nil
}
