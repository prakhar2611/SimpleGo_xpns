package Models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId    uint32 `gorm:"primaryKey;autoIncrement:true"`
	FirstName string
	LastName  string
	Email     string
	Address   string
}

type UserToken struct {
	gorm.Model
	User  User `gorm:"foreignKey:UserId"`
	Token string
}
