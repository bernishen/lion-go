package sessionservice

import (
	"bytes"
	"encoding/json"
	"github.com/Berni-Shen/lion-go/utils/exception"
	"strings"
	"time"

	"github.com/Berni-Shen/lion-go/oauth2/dal/domain"
	"github.com/Berni-Shen/lion-go/oauth2/service/redisservice"
)

const (
	signLogTag = "signuser:"
	globalTag  = "Global:"
)

type Session struct {
	UserID string
	Roles  map[string][]domain.Role
}

// NewGlobal : Create a global session or update the global session the expire time.
func NewGlobal(userID string, accessToken string, roles *[]domain.Role) (bool, *exception.Exception) {
	key := getGlobalKey(accessToken)
	exits, _ := redisservice.Exits(key)

	if exits {
		return redisservice.RefreshExpire(key, time.Minute*15)
	}

	s := Session{
		UserID: userID,
		Roles:  make(map[string][]domain.Role),
	}
	for _, r := range *roles {
		rs, ok := s.Roles[r.SystemID]
		if !ok {
			s.Roles[r.SystemID] = []domain.Role{r}
		} else {
			s.Roles[r.SystemID] = append(rs, r)
		}
	}
	sJson, err := json.Marshal(s)
	if err != nil {
		return false, exception.NewException(exception.Error, 1001, "Had a error occurred marshal the roles infomation.["+err.Error()+"]")
	}
	var builder strings.Builder
	builder.Write(sJson)
	redisservice.Set(signLogTag+userID, accessToken, time.Minute*15)

	return redisservice.Set(key, builder.String(), time.Minute*15)
}

// VerifyGlobal : Verified the user is login in global, and return this user's role.
func VerifyGlobal(accessToken string) (*Session, *exception.Exception) {
	key := getGlobalKey(accessToken)
	exits, ex := redisservice.Exits(key)
	if ex != nil {
		return nil, ex.ResetCode(1001)
	}
	if !exits {
		return nil, exception.NewException(exception.Warning, 1002, "The user is not login or login is to longer.")
	}
	rString, ex1 := redisservice.Get(key)
	if ex1 != nil {
		return nil, ex1.ResetCode(1003)
	}
	var buf bytes.Buffer
	buf.WriteString(rString)
	var session Session
	json.Unmarshal(buf.Bytes(), &session)

	return &session, nil
}

// getGlobalKey : Geting a key of the global session.
func getGlobalKey(accessToken string) string {
	var buider strings.Builder
	buider.WriteString(globalTag)
	buider.WriteString(accessToken)

	return buider.String()
}
