package rsaservice

import (
	"fmt"
	"github.com/bernishen/lion-go/oauth2/dal/dbpool"
	"github.com/bernishen/lion-go/oauth2/dal/domain"
	"github.com/bernishen/lion-go/oauth2/service/redisservice"
	"github.com/bernishen/exception"
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
