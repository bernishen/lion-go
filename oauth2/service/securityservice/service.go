package securityservice

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/bernishen/lion-go/oauth2/service/redisservice"
	"github.com/bernishen/lion-go/utils/exception"
	"strings"
	"time"
)

type iservice interface {
	PrivateKey() (string, *exception.Exception)
	PublicKey() (string, *exception.Exception)
	Decrypt(data string) (string, *exception.Exception)
}

type key struct {
	PrivateKey string `json:"private_key"`
	Version    string `json:"version""`
}

func instance(version string, privateKey string) iservice {
	var s iservice
	var ex *exception.Exception
	switch {
	case version == "1.0" || version == "1.1":
		s, ex = initRSA(version, privateKey)
	default:
		s, ex = initRSA(version, privateKey)
	}
	if ex != nil {
		return nil
	}
	return s
}

// FindPrivateKey Find the private key by 'clientID'.
// Used to obtain the cryptographic key.
func FindPrivateKey(clientID string) (string, *exception.Exception) {
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

func DecryptData(clientID string, data string) (string, *exception.Exception) {
	dBytes, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		return "", exception.NewException(exception.Error, 1001, "The data is not base64.-->"+err.Error())
	}
	var builder strings.Builder
	builder.Write(dBytes)

	keyText, ex := FindPrivateKey(clientID)
	if ex != nil {
		return "", ex
	}

	var buffer bytes.Buffer
	_, err = buffer.WriteString(keyText)
	if err != nil {
		return "", exception.NewException(exception.Error, 1002, "The read key failed.-->"+err.Error())
	}

	var k key
	json.Unmarshal(buffer.Bytes(), &k)

	s := instance(k.Version, k.PrivateKey)
	plaintext, ex := s.Decrypt(builder.String())
	if ex != nil {
		return "", ex
	}
	return plaintext, nil
}

func CreateKey(clientID string, version string) (string, *exception.Exception) {
	s := instance(version, "")

	privateKey, ex := s.PrivateKey()
	if ex != nil {
		return "", ex
	}

	itemBytes, err := json.Marshal(key{
		privateKey,
		version,
	})
	if err != nil {
		return "", exception.NewException(exception.Error, 1001, "The can't marshal.-->"+err.Error())
	}

	var builder strings.Builder
	builder.Write(itemBytes)
	_, ex = redisservice.Set(clientID, builder.String(), time.Minute*5)
	if ex != nil {
		return "", ex
	}

	publicKey, ex := s.PublicKey()
	if ex != nil {
		return "", ex
	}

	var buffer bytes.Buffer
	buffer.WriteString(publicKey)
	publicKey = base64.RawURLEncoding.EncodeToString(buffer.Bytes())

	return publicKey, nil
}
