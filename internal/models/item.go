package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Item struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name        string         `gorm:"not null;index" json:"name"`
	Description *string        `json:"description,omitempty"`
	Quantity    int            `gorm:"not null;default:0" json:"quantity"`
}

func (Item) TableName() string {
	return "items"
}

func (i *Item) BeforeCreate(_ *gorm.DB) error {
	i.Name = strings.TrimSpace(i.Name)
	return nil
}

func (i *Item) BeforeUpdate(_ *gorm.DB) error {
	i.Name = strings.TrimSpace(i.Name)
	return nil
}
