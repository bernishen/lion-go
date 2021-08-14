package sessionservice

import (
	"github.com/bernishen/lion-go/utils/exception"
	"strings"
	"time"

	"github.com/bernishen/lion-go/oauth2/dal/domain"
	"github.com/bernishen/lion-go/oauth2/service/redisservice"
	"github.com/bernishen/lion-go/oauth2/service/tokenservice"
)

// GetToken : Find the token from the system session.
func GetToken(systemID string, accessToken string) (string, *exception.Exception) {
	key := getSystemKey(systemID, accessToken)

	return redisservice.Get(key)
}

// NewSystem : Create a system session or update the system session the expire time.
func NewSystem(systemID string, accessToken string, roles *[]domain.Role) (string, *exception.Exception) {
	key := getSystemKey(systemID, accessToken)
	exits, _ := redisservice.Exits(key)
	if exits {
		redisservice.RefreshExpire(key, time.Minute*15)
		return redisservice.Get(key)
	}
	token, ex := tokenservice.SignToken(systemID, accessToken, roles)
	if ex != nil {
		return "", ex.ResetCode(1001)
	}
	ok, ex1 := redisservice.Set(key, token, time.Minute*15)
	if ex1 != nil {
		return "", ex1.ResetCode(1002)
	}
	if !ok {
		return "", exception.NewException(exception.Error, 1003, "The session save to cache is failed.")
	}

	return token, nil
}

// VerifySystem : Verified the user is login in this system.
func VerifySystem(systemID string, accessToken string) (string, *exception.Exception) {
	key := getSystemKey(systemID, accessToken)
	exist, ex := redisservice.Exits(key)
	if ex != nil || !exist {
		return "", ex
	}
	token, ex := redisservice.Get(key)
	if ex != nil {
		return "", ex
	}
	return token, nil
}

// getSystemKey : Geting a key of the global session.
func getSystemKey(systemID string, accessToken string) string {
	var buider strings.Builder
	buider.WriteString("System[")
	buider.WriteString(systemID)
	buider.WriteString("]:")
	buider.WriteString(accessToken)

	return buider.String()
}
