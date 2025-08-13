package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title   string `binding:"required"`
	Content string `binding:"required"`
	Preview string `binding:"required"`
	// 这边把likes的数据删掉，防止redis与数据库这项数据不同步
}
