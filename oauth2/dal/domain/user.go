package domain

// User :
type User struct {
	BaseModel
	UserID          string `gorm:"type:varchar(50);comment:'用户标识'"`
	Password        string `gorm:"type:varchar(100);comment:'密码'"`
	PasswordVersion int    `gorm:"not_null:size:2;comment:'密码版本'"`
	UserState       int    `gorm:"not_null;size:1;comment:'用户状态^0：停用，1：启用。'"`
	UserName        string `gorm:"type:varchar(100);comment:'用户姓名'"`
	Email           string `gorm:"type:varchar(100);comment:'电子邮箱'"`
	Phone           string `gorm:"type:varchar(20);comment:'电话'"`
}
