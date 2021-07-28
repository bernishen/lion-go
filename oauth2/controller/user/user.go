package user

import (
	"github.com/Berni-Shen/lion-go/oauth2/service/accountservice"
	"github.com/Berni-Shen/lion-go/oauth2/service/accountservice/domain"
	"github.com/Berni-Shen/lion-go/utils/exception"
	"github.com/Berni-Shen/lion-go/utils/router"
)

func init() {
	u := router.InitController("/oauth2/account/:userid/:pwd").
		Get(signIn, "userid", "pwd").
		Post(signUp).
		Delete(signOut).
		Delete(accountDelete, "userid", "pwd")
	router.Default().Register(u)
}

type user struct {
	ClientID   string `json:"client_id"`
	UserID     string `json:"user_id"`
	Password   string `json:"password"`
	UserName   string `json:"user_name"`
	UserCardID string `json:"user_card_id"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}

func (u *user) toSignUp() *domain.SignUpUser {
	user := domain.SignUpUser{
		ClientID:   u.ClientID,
		UserID:     u.UserID,
		Password:   u.Password,
		UserName:   u.UserName,
		UserCardID: u.UserCardID,
		Email:      u.Email,
		Phone:      u.Phone,
	}
	return &user
}

func signUp(u user) (*interface{}, *exception.Exception) {
	user := u.toSignUp()
	ok, ex := accountservice.SignUp(user)
	if ex != nil {
		return nil, exception.NewException(exception.Error, 1001, "Sign up failed.[msg:"+ex.Message+"]")
	}
	if !ok {
		return nil, exception.NewException(exception.Error, 1001, "Sign up failed.[msg:Inner error]")
	}
	token, ex := accountservice.SignIn(u.ClientID, u.UserID, u.Password)
	if ex != nil {
		return nil, exception.NewException(exception.Warning, 1002, "Sign up success, but sign in failed.[msg:"+ex.Message+"]")
	}

	var ret interface{} = struct {
		Token string `json:"token"`
	}{
		token,
	}
	return &ret, nil
}

func signIn(userid string, pwd string, ctx *router.Context) (*interface{}, *exception.Exception) {
	cid := ctx.Request.Header.Get("cid")
	token, ex := accountservice.SignIn(cid, userid, pwd)
	if ex != nil {
		return nil, exception.NewException(exception.Error, 1001, ex.Message)
	}

	var ret interface{} = struct {
		Token string `json:"token"`
	}{
		token,
	}
	return &ret, nil
}

func signOut(ctx *router.Context) (*interface{}, *exception.Exception) {
	token := ctx.Token
	var ret interface{} = struct {
		Token string `json:"token"`
	}{
		token,
	}
	ex := accountservice.SignOut(token)
	if ex != nil {
		return &ret, ex
	}
	return &ret, nil
}

func accountDelete(userid string, pwd string, ctx *router.Context) (*interface{}, *exception.Exception) {
	var ret interface{} = "This is accountDelete."
	return &ret, nil
}
