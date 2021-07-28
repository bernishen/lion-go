package domain

// UserNRole : Connection the users and the roles.
type UserNRole struct {
	BaseModel
	UserID string `gorm:"type:char(36);not_null"`
	RoleID string `gorm:"type:char(36);not_null"`
}
