package task

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Task   string `json:"task"`
	IsDone bool   `json:"is_done" gorm:"not null"`
	UserID uint32 `json:"user_id" gorm:"not null;index"`
	User   User   `gorm:"constraint:OnDelete:CASCADE;"` // Связь с каскадным удалением
}

type User struct {
	ID uint32 `gorm:"primaryKey"`
}
