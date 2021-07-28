package redisservice

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/Berni-Shen/lion-go/oauth2/common/dao/domain"
	"github.com/Berni-Shen/lion-go/oauth2/service/redisservice"
)

func TestSetKey(t *testing.T) {
	// for i := 0; i < 5; i++ {
	// 	redisservice.Set(i, i, time.Second*30)
	// }
}

func TestGet(t *testing.T) {
	v, err := redisservice.Get("Global:123")
	var buffer bytes.Buffer
	buffer.WriteString(v)
	var m map[string][]domain.Role
	err1 := json.Unmarshal(buffer.Bytes(), &m)

	t.Log(m)
	t.Log(err)
	t.Log(err1)
}
