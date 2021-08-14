package dao

import (
	"github.com/bernishen/lion-go/oauth2/common/dao/dbskim"
	"testing"
)

func TestInitTables(t *testing.T) {
	_, ex := dbskim.InitTables()
	if ex !=nil{
		t.Log(ex.Message)
	}
}
