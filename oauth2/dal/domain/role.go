package domain

// Role :
type Role struct {
	BaseModel
	Name        string `gorm:"type:varchar(50);not_null"`
	Description string `gorm:"type:varchar(200)"`
	Enabled     bool   `gorm:"not_null"`
	SystemID    string `gorm:"type:varchar(50);not_null"`
}
