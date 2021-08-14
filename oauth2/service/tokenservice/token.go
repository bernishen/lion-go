package tokenservice

import (
	"github.com/bernishen/lion-go/oauth2/dal/domain"
	"github.com/bernishen/lion-go/utils/exception"
)

func SignToken(accessToken string, systemID string, roles *[]domain.Role) (string, *exception.Exception) {

	return "", nil
}
