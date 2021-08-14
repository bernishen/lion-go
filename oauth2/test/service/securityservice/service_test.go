package securityservice

import (
	"bytes"
	"encoding/base64"
	"github.com/bernishen/lion-go/oauth2/service/securityservice"
	"testing"
)

func TestBase64(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("abcd")
	cipherText1 := base64.StdEncoding.EncodeToString(buffer.Bytes())
	t.Logf("ciphertext1:" + cipherText1)
	cipherText2 := base64.URLEncoding.EncodeToString(buffer.Bytes())
	t.Logf("ciphertext2:" + cipherText2)
	cipherText3 := base64.RawStdEncoding.EncodeToString(buffer.Bytes())
	t.Logf("ciphertext3:" + cipherText3)
	cipherText4 := base64.RawURLEncoding.EncodeToString(buffer.Bytes())
	t.Logf("ciphertext4:" + cipherText4)
}

func TestCreateKey(t *testing.T) {
	s, ex := securityservice.CreateKey("1", "1.0")
	if ex != nil {
		t.Log(ex.Message)
	}
	var buffer bytes.Buffer
	buffer.WriteString(s)
	cipherText := base64.StdEncoding.EncodeToString(buffer.Bytes())
	t.Logf("ciphertext:" + cipherText)
}
