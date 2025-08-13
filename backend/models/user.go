package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	// 保证用户名不重复
	Username string `gorm:"unique"`
	// 密码
	Password string
}
