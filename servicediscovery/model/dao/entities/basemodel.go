package entities

import (
	"time"
)

// BaseModel : this is basic fields in database tables.
type BaseModel struct {
	ID        string `gorm:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
