package db

import "gorm.io/gorm"

type Tabler interface {
	TableName() string
}

func (User) TableName() string {
	return "User"
}

type User struct {
	gorm.Model
	Account  string `gorm:"type:varchar(32);not null;unique_index"`
	Password string `gorm:"type:varchar(32);not null;default:'123456'"`
}
