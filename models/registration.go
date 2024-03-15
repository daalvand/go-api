package models

import "gorm.io/gorm"

type Registration struct {
	gorm.Model
	EventID uint
	UserID  uint
}
