package domain

import (
	"github.com/google/uuid"
	"time"
)

// BaseModel is
type BaseModel struct {
	ID        string    `gorm:"type:char(36);primary_key;comment:'标识'"`
	CreatedAt time.Time `gorm:"comment:'创建时间'"`
	UpdatedAt time.Time `gorm:"comment:'更新时间'"`
	DeletedAt *time.Time `gorm:"comment:'删除时间'"`
}

// New
func NewBase() *BaseModel {
	return &BaseModel{
		ID:        uuid.New().String(),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: nil,
	}
}