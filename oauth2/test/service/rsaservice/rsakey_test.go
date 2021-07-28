package rsaservice

import (
	"github.com/Berni-Shen/lion-go/oauth2/service/rsaservice"
	"testing"
)

func TestCreateKeys(t *testing.T) {
	t.Log(rsaservice.CreateKey("aaaabbbccc"))
}
