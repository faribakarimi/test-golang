package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID			int		`gorm:"primary_key;auto_increment" json:"id"`
	Username  	string	`gorm:"size:255;not null" json:"username"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Name     	string	`gorm:"size:255;not null" json:"name"`
	Family  	string	`gorm:"size:255;not null;" json:"family"`
	Gender		string	`gorm:"size:10;not null" json:"gender"`
	Age			uint16	`gorm:"not null" json:"age"`
	Balance		uint64	`gorm:"not null" json:"balance"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (user *User) BeforeSave() error {
	hashedPassword, err := Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}