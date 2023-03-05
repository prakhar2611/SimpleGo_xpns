package Models

import (
	"time"

	"gorm.io/gorm"
)

type UserToken struct {
	gorm.Model
	ID           string `gorm:"primaryKey;unique"`
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time
}

type User struct {
	gorm.Model
	ID        string `gorm:"nullable"`
	Email     string `gorm:"index;unique"`
	Name      string
	Picture   string
	Locale    string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
