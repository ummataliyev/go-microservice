package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Email          string         `gorm:"uniqueIndex;not null" json:"email"`
	HashedPassword string         `gorm:"not null" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	return nil
}

func (u *User) BeforeUpdate(_ *gorm.DB) error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	return nil
}
