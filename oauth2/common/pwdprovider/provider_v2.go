package pwdprovider

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"github.com/Berni-Shen/lion-go/utils/exception"
)

type providerV2 struct {
}

// Compute : the plaintext compute to cryphertext.
func (p *providerV2) Compute(plainText string) (string, *exception.Exception) {
	s := sha256.New()
	var buffer bytes.Buffer
	buffer.WriteString(plainText)
	_, err := s.Write(buffer.Bytes())
	if err != nil {
		return "", exception.NewException(exception.Error, 2001, "Readed 'plainText' error.["+err.Error()+"]")
	}
	cryphers := s.Sum(nil)

	return base64.StdEncoding.EncodeToString(cryphers), nil
}
