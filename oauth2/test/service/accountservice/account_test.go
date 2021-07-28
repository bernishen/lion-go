package accountservice

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/Berni-Shen/lion-go/oauth2/service/rsaservice"
	"github.com/google/uuid"
	"testing"

	"github.com/Berni-Shen/lion-go/oauth2/service/accountservice"
	"github.com/Berni-Shen/lion-go/oauth2/service/accountservice/domain"
)

func TestSignUp(t *testing.T) {
	u := domain.UserInfo{
		ClientID:   "shenglsekdhlwnga",
		UserID:     "admin",
		UserCardID: "530103000000000001",
		UserName:   "Berni Shen",
		Password:   "123",
		Phone:      "15966666666",
		Email:      "test@lion-go.net",
	}

	u.Password = encryptPassword(u.ClientID, u.Password)

	b, ex := accountservice.SignUp(&u)
	if !b {
		t.Log(ex.Message)
	}
}

func TestDelete(t *testing.T) {
	d, ex := accountservice.Delete("admin")
	if !d {
		t.Log(ex.Message)
	}
}

func TestSingIn(t *testing.T) {
	cid := uuid.New().String()
	token, ex := accountservice.SignIn(cid, "Admin", encryptPassword(cid, "123"))
	if ex != nil {
		t.Log("Error-->" + ex.Message)
	}
	t.Log("Success-->" + token)
}

func encryptPassword(clientID string, pwd string) string {
	rsaKey, _ := rsaservice.CreateKey(clientID)

	block, _ := pem.Decode([]byte(rsaKey))
	if block == nil {
		panic("failed to parse PEM block containing the public key")
	}
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}

	var buffer bytes.Buffer
	buffer.WriteString(pwd)
	pwdNew, err := rsa.EncryptPKCS1v15(rand.Reader, pub, buffer.Bytes())

	buffer.Reset()
	buffer.Write(pwdNew)

	return buffer.String()
}
