package pwdprovider

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"github.com/bernishen/exception"
)

type providerV1 struct {
}

// Compute : the plaintext compute to cryphertext.
func (p *providerV1) Compute(plainText string) (string, *exception.Exception) {
	s := md5.New()
	var buffer bytes.Buffer
	buffer.WriteString(plainText)
	_, err := s.Write(buffer.Bytes())
	if err != nil {
		return "", exception.NewException(exception.Error, 2001, "Readed 'plainText' error.["+err.Error()+"]")
	}
	cryphers := s.Sum(nil)

	return base64.StdEncoding.EncodeToString(cryphers), nil
}
