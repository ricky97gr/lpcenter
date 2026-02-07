package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID      string    `json:"uuid" gorm:"column:uuid;type:varchar(36);uniqueIndex"`
	Username  string    `json:"username" gorm:"column:username;type:varchar(50);uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"column:email;type:varchar(100);uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"column:password;type:varchar(255);not null"`
	Status    string    `json:"status" gorm:"column:status;type:varchar(20);default:'active'"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusDisabled = "disabled"
)

func (User) TableName() string {
	return "users"
}
