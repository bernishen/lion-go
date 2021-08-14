package securityservice

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/bernishen/lion-go/utils/exception"
	"strings"
)

type irsaitem interface {
	privateKey() ([]byte, *exception.Exception)
	publicKey() ([]byte, *exception.Exception)
	readPrivateKey(key []byte) (*rsa.PrivateKey, *exception.Exception)
	decrypt(data []byte) ([]byte, *exception.Exception)
}

type rsa_base struct {
	key  *rsa.PrivateKey
	item irsaitem
}

func initRSA(version string, privateKey string) (iservice, *exception.Exception) {
	version = strings.Trim(version, " ")
	version = strings.ToLower(version)

	k := &rsa_base{}
	switch {
	case version == "1.0":
		item := &rsa_pkcs1{
			base: k,
		}
		ex := k.init(item, privateKey, 1024)
		if ex != nil {
			return nil, ex
		}
	case version == "1.1" || version == "lastest":
		item := &rsa_pkix{
			base: k,
		}
		ex := k.init(item, privateKey, 1024)
		if ex != nil {
			return nil, ex
		}
	default:
		return nil, exception.NewException(exception.Error, 1002, "Found not security service of version is '"+version+"'")
	}
	return k, nil
}

func (k *rsa_base) init(item irsaitem, privateKey string, size int) *exception.Exception {
	var key *rsa.PrivateKey
	var err error
	if privateKey == "" {
		key, err = rsa.GenerateKey(rand.Reader, size)
		if err != nil {
			return exception.NewException(exception.Error, 1001, "Has an exception current.-->"+err.Error())
		}
	} else {
		var buffer bytes.Buffer
		buffer.WriteString(privateKey)
		block, _ := pem.Decode(buffer.Bytes())
		if block == nil {
			return exception.NewException(exception.Error, 1001, "An exception occurred reading the 'PEM' key, the block can't null.")
		}
		var ex *exception.Exception
		key, ex = item.readPrivateKey(block.Bytes)
		if ex != nil {
			return ex
		}
	}
	k.key = key
	k.item = item
	return nil
}

func (k *rsa_base) PrivateKey() (string, *exception.Exception) {
	if k.item == nil || k.key == nil {
		return "", exception.NewException(exception.Error, 1001, "Initialization rsa failed.")
	}

	key, ex := k.item.privateKey()
	if ex != nil {
		return "", ex
	}

	builder := &strings.Builder{}
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: key,
	}
	pem.Encode(builder, block)

	return builder.String(), nil
}

func (k *rsa_base) PublicKey() (string, *exception.Exception) {
	if k.item == nil || k.key == nil {
		return "", exception.NewException(exception.Error, 1001, "Initialization rsa failed.")
	}

	key, ex := k.item.publicKey()
	if ex != nil {
		return "", ex
	}

	builder := &strings.Builder{}
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: key,
	}
	pem.Encode(builder, block)

	return builder.String(), nil
}

func (k *rsa_base) Decrypt(data string) (string, *exception.Exception) {
	var buffer bytes.Buffer
	buffer.WriteString(data)
	plainByte, ex := k.item.decrypt(buffer.Bytes())
	if ex != nil {
		return "", ex
	}

	var builder strings.Builder
	for i := 0; i < len(plainByte); i++ {
		builder.WriteByte(plainByte[i])
	}

	return builder.String(), nil
}

type rsa_pkcs1 struct {
	base *rsa_base
}

func (k *rsa_pkcs1) privateKey() ([]byte, *exception.Exception) {
	if k.base.key == nil {
		return nil, exception.NewException(exception.Error, 1001, "The key is not null.")
	}
	return x509.MarshalPKCS1PrivateKey(k.base.key), nil
}

func (k *rsa_pkcs1) publicKey() ([]byte, *exception.Exception) {
	if k.base.key == nil {
		return nil, exception.NewException(exception.Error, 1001, "The key is not null.")
	}
	pub := k.base.key.PublicKey
	return x509.MarshalPKCS1PublicKey(&pub), nil
}

func (k *rsa_pkcs1) readPrivateKey(key []byte) (*rsa.PrivateKey, *exception.Exception) {
	ret, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, exception.NewException(exception.Error, 1001, "Read 'private key' error.-->"+err.Error())
	}
	return ret, nil
}

func (k *rsa_pkcs1) decrypt(data []byte) ([]byte, *exception.Exception) {
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, k.base.key, data)
	if err != nil {
		return nil, exception.NewException(exception.Error, 1001, "Decrypt the text failed.-->"+err.Error())
	}
	return plaintext, nil
}

type rsa_pkix struct {
	base *rsa_base
}

func (k *rsa_pkix) privateKey() ([]byte, *exception.Exception) {
	if k.base.key == nil {
		return nil, exception.NewException(exception.Error, 1001, "The key is not null.")
	}
	return x509.MarshalPKCS1PrivateKey(k.base.key), nil
}

func (k *rsa_pkix) publicKey() ([]byte, *exception.Exception) {
	if k.base.key == nil {
		return nil, exception.NewException(exception.Error, 1001, "The key is not null.")
	}
	pub := k.base.key.PublicKey
	pubByte, err := x509.MarshalPKIXPublicKey(&pub)
	if err != nil {
		return nil, exception.NewException(exception.Error, 1002, "Marshal public key failed.-->"+err.Error())
	}

	return pubByte, nil
}

func (k *rsa_pkix) readPrivateKey(key []byte) (*rsa.PrivateKey, *exception.Exception) {
	ret, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, exception.NewException(exception.Error, 1001, "Read 'private key' error.-->"+err.Error())
	}
	return ret, nil
}

func (k *rsa_pkix) decrypt(data []byte) ([]byte, *exception.Exception)  {
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, k.base.key, data)
	if err != nil {
		return nil, exception.NewException(exception.Error, 1001, "Decrypt the text failed.-->"+err.Error())
	}
	return plaintext, nil
}
