package tokenservice

import (
	"github.com/Berni-Shen/lion-go/oauth2/dal/domain"
	"github.com/Berni-Shen/lion-go/utils/exception"
)

func SignToken(accessToken string, systemID string, roles *[]domain.Role) (string, *exception.Exception) {

	return "", nil
}
