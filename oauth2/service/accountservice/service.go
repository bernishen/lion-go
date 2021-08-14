package accountservice

import (
	"github.com/bernishen/lion-go/oauth2/common/pwdprovider"
	"github.com/bernishen/lion-go/oauth2/dal/dbpool"
	"github.com/bernishen/lion-go/oauth2/dal/domain"
	pd "github.com/bernishen/lion-go/oauth2/service/accountservice/domain"
	"github.com/bernishen/lion-go/oauth2/service/pwdservice"
	"github.com/bernishen/lion-go/oauth2/service/redisservice"
	"github.com/bernishen/lion-go/oauth2/service/sessionservice"
	"github.com/bernishen/lion-go/oauth2/service/tokenservice"
	"github.com/bernishen/lion-go/utils/exception"
)

const (
	signLogTag = "signuser"
)

// SignUp is user created a account for sign in this system.
func SignUp(user *pd.SignUpUser) (bool, *exception.Exception) {
	db, ex := dbpool.Take()
	if ex != nil {
		return false, exception.NewException(exception.Error, 1001, ex.Message)
	}
	defer dbpool.Put(db)

	var users []domain.User
	db.Find(&users, "user_id = ?", user.UserID)
	if len(users) > 0 {
		return false, exception.NewException(exception.Error, 1002, "This 'UseID' is already in use. Please changed it and try again.")
	}
	u := &domain.User{
		BaseModel:       *domain.NewBase(),
		UserID:          user.UserID,
		UserName:        user.UserName,
		Phone:           user.Phone,
		Email:           user.Email,
		PasswordVersion: pwdprovider.NowPasswordVersion(),
		UserState:       1,
	}

	u.Password, ex = pwdservice.ConvertPwd(user.ClientID, user.Password, u.PasswordVersion)
	if ex != nil {
		return false, exception.NewException(exception.Error, 1003, "System Error, password converted failed.-->"+ex.Message)
	}

	db.Create(&u)
	var u2 domain.User
	db.First(&u2)
	if &u2 == nil {
		return false, exception.NewException(exception.Error, 1004, "An exception has occaurred. Please try again.")
	}
	return true, nil
}

// Delete is user delete the account.
func Delete(password string, clientID string, accessToken string) (bool, *exception.Exception) {
	s, ex := sessionservice.VerifyGlobal(accessToken)
	if ex != nil {
		return false, ex
	}

	u, ex := verifyUser(clientID, s.UserID, password)
	if ex != nil {
		return false, ex
	}

	db, ex := dbpool.Take()
	if ex != nil {
		return false, exception.NewException(exception.Error, 1001, ex.Message)
	}
	defer dbpool.Put(db)
	db.Delete(&u)
	return true, nil
}

// SignIn the system by the UserID and Password, Authorize the user is legal.
// Return the access token.
func SignIn(clientID string, userID string, password string) (string, *exception.Exception) {
	exists, ex := redisservice.Exits(signLogTag + userID)
	if exists {
		at, ex := redisservice.Get(signLogTag + userID)
		if ex != nil {
			_, ex := sessionservice.VerifyGlobal(at)
			if ex == nil {
				return at, nil
			}
		}
	}

	_, ex = verifyUser(clientID, userID, password)
	if ex != nil {
		return "", ex
	}

	roles, ex := findRoles(userID)
	if ex != nil {
		return "", ex.ResetCode(1005)
	}

	accessToken, ex := tokenservice.SignAccessToken(clientID)
	if ex != nil {
		return "", ex.ResetCode(1006)
	}

	_, ex = sessionservice.NewGlobal(userID, accessToken, roles)
	if ex != nil {
		return "", ex.ResetCode(1007)
	}
	_, ex = redisservice.Set(signLogTag+userID, accessToken, -1)

	return accessToken, nil
}

func SignOut(accessToken string) *exception.Exception {
	_, ex := sessionservice.VerifyGlobal(accessToken)
	if ex != nil {
		return ex
	}
	//ex = sessionservice.RemoveGlobal(accessToken)
	if ex != nil {
		return ex
	}
	return nil
}

func SignVerify(clientID string, accessToken string, password string) (bool, *exception.Exception) {
	s, ex := sessionservice.VerifyGlobal(accessToken)
	if ex != nil {
		return false, ex
	}

	_, ex = verifyUser(clientID, s.UserID, password)
	if ex != nil {
		return false, ex
	}
	return true, nil
}

func verifyUser(clientID string, userid string, password string) (*domain.User, *exception.Exception) {
	db, ex := dbpool.Take()
	if ex != nil {
		return nil, exception.NewException(exception.Error, 1001, ex.Message)
	}
	defer dbpool.Put(db)
	var u domain.User
	db.First(&u, "user_id = ?", userid)
	if &u == nil {
		return nil, exception.NewException(exception.Error, 1002, "Found not the user.")
	}
	pwd, ex := pwdservice.ConvertPwd(clientID, password, u.PasswordVersion)
	if ex != nil {
		return nil, ex.ResetCode(1003)
	}
	if pwd != u.Password {
		return nil, exception.NewException(exception.Error, 1004, "Verify failed, becuse the UserID or Password is error.")
	}

	return &u, nil
}

// RegisterSystem : Registe this user in the system.
func RegisterSystem(systemID string, accessToken string) (string, *exception.Exception) {
	//global, ex := sessionservice.VerifyGlobal(accessToken)
	//if ex != nil {
	//	return "", ex.ResetCode(1001)
	//}

	//conn, ok := broadcast.ConnPool[systemID]
	//if ok {
	//
	//}
	//conn.WriteMessage()

	system, _ := sessionservice.VerifySystem(systemID, accessToken)
	if system == "" {
		return sessionservice.GetToken(systemID, accessToken)
	}
	r := make((map[string][]domain.Role), 1)
	roles := &r
	systemRoles, ok := (*roles)[systemID]
	if !ok {
		return "", exception.NewException(exception.Error, 1002, "This user had not authority in the system.")
	}
	return sessionservice.NewSystem(systemID, accessToken, &systemRoles)
}

func findRoles(userID string) (*[]domain.Role, *exception.Exception) {
	db, ex := dbpool.Take()
	if ex != nil {
		return nil, ex
	}
	defer dbpool.Put(db)
	var unr []domain.UserNRole
	db.Find(&unr, "user_id = ?", userID)
	roles := make([]domain.Role, len(unr))
	for i := 0; i < len(unr); i++ {
		item := unr[i]
		db.First(&roles[i], "ID = ?", item.RoleID)
	}

	return &roles, nil
}
