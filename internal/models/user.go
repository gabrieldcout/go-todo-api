package models

type User struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Email        string `json:"email" gorm:"unique"`
	PasswordHash string `json:"-"`
}
