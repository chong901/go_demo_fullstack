package filter

import "github.com/jinzhu/gorm"

type Filter interface {
	SetCriteria(db *gorm.DB) *gorm.DB
}
