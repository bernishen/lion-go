package dbskim

import (
	"github.com/Berni-Shen/lion-go/oauth2/common/dao/dbpool"
	"github.com/Berni-Shen/lion-go/oauth2/common/dao/domain"
	"github.com/Berni-Shen/lion-go/utils/exception"
)

func InitTables() (bool, *exception.Exception) {
	db, ex := dbpool.Take()
	if ex != nil {
		ex.NewCode(1001)
	}
	db.CreateTable(domain.User{})
	db.CreateTable(domain.Role{})
	db.CreateTable(domain.UserNRole{})
	db.CreateTable(domain.RSAKey{})
	return true, nil
}
