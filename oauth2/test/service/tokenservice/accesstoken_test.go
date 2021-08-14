package tokenservice

import (
	"testing"

	"github.com/bernishen/lion-go/oauth2/common/dao/domain"
	"github.com/bernishen/lion-go/oauth2/service/sessionservice"
)

func TestSign(t *testing.T) {
	roles := make([]domain.Role, 3)
	roles[0].ID = "1"
	roles[0].Name = "1"
	roles[0].SystemID = "1"
	roles[1].ID = "2"
	roles[1].Name = "2"
	roles[1].SystemID = "1"
	roles[2].ID = "3"
	roles[2].Name = "3"
	roles[2].SystemID = "3"

	t.Log(sessionservice.NewGlobal("123", &roles))
}
