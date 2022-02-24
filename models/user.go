package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string
}

func CheckAuth(username, password string) (bool, error) {
	var user User
	err := db.Select("id").Where(User{Username: username, Password: password}).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if user.ID > 0 {
		return true, nil
	}

	return false, nil
}