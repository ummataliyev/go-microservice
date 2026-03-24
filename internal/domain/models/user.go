package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// User represents the users table in the database.
type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Email          string         `gorm:"uniqueIndex;not null" json:"email"`
	HashedPassword string         `gorm:"not null" json:"-"`
}

// TableName returns the table name for the User model.
func (User) TableName() string {
	return "users"
}

// BeforeCreate normalises the email to lowercase before inserting.
func (u *User) BeforeCreate(_ *gorm.DB) error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	return nil
}

// BeforeUpdate normalises the email to lowercase before updating.
func (u *User) BeforeUpdate(_ *gorm.DB) error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	return nil
}
