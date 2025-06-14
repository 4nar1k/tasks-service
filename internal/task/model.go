package task

type Task struct {
	ID     uint32 `gorm:"primaryKey;autoIncrement" json:"id"`
	Title  string `gorm:"type:text;not null" json:"title"`
	IsDone bool   `gorm:"not null;default:false" json:"is_done"`
	UserID uint32 `gorm:"not null;index" json:"user_id"`
	User   User   `gorm:"constraint:OnDelete:CASCADE;"`
}

type User struct {
	ID uint32 `gorm:"primaryKey" json:"id"`
}
