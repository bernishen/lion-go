package rsaservice

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/Berni-Shen/lion-go/oauth2/service/redisservice"
	"github.com/Berni-Shen/lion-go/utils/exception"
	"strings"
	"time"

	"github.com/Berni-Shen/lion-go/oauth2/dal/dbpool"
	"github.com/Berni-Shen/lion-go/oauth2/dal/domain"
)

func FindPublicKey(scope string, version int) (string, error) {
	db, ex := dbpool.Take()
	if ex != nil {
		return "", fmt.Errorf("Can't find a public key, because : %s", ex.Message)
	}
	var key domain.RSAKey
	db.Where("scope = ? and version = ?", scope, version).First(&key)
	if &key == nil {
		return "", fmt.Errorf("Found not the public key, [srope:%s,version:%d]", scope, version)
	}
	return key.PublicKey, nil
}

func FindPrivateKey(scope string, version int) (string, error) {
	db, ex := dbpool.Take()
	if ex != nil {
		return "", fmt.Errorf("Can't find a private key, because : %s", ex.Message)
	}
	var key domain.RSAKey
	db.Where("scope = ? and version = ?", scope, version).First(&key)
	if &key == nil {
		return "", fmt.Errorf("Found not the private key, [srope:%s,version:%d]", scope, version)
	}
	return key.PrivateKey, nil
}

func CreateKey(clientID string) (string, *exception.Exception) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", exception.NewException(exception.Error, 1001, "Has an exception current.-->"+err.Error())
	}

	priBytes := x509.MarshalPKCS1PrivateKey(key)
	builder := &strings.Builder{}
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: priBytes,
	}
	err = pem.Encode(builder, block)
	if err != nil {
		return "", exception.NewException(exception.Error, 1002, "Has an exception occurred during production of the PEM files.==>"+err.Error())
	}

	_, ex := redisservice.Set(clientID, builder.String(), time.Minute*5)
	builder.Reset()
	if ex != nil {
		return "", ex.ResetCode(1003)
	}

	pubBytes := x509.MarshalPKCS1PublicKey(&key.PublicKey)
	block = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubBytes,
	}
	pem.Encode(builder, block)

	return builder.String(), nil
}

// FindKeyByClient Find the private key by 'clientID'.
// Used to obtain the cryptographic key.
func FindKeyByClient(clientID string) (string, *exception.Exception) {
	exits, ex := redisservice.Exits(clientID)
	if ex != nil {
		return "", ex.ResetCode(1001)
	}

	if !exits {
		return "", exception.NewException(exception.Error, 1002, "This client security session was not found.")
	}

	key, ex1 := redisservice.Get(clientID)
	if ex1 != nil {
		return "", ex1.ResetCode(1001)
	}

	return key, nil
}
