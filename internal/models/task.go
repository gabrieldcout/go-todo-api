package models

type Task struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`

	UserID uint `json:"-"`
}
