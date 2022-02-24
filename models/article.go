package models

import (
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title   string
	Content string
	UserID  uint
}
