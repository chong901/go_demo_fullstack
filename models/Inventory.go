package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Inventory struct {
	Id      uint               `gorm:"primary_key" json:"id" form:"id"`
	Level   int                `json:"level"`
	User    string             `json:"user"`
	SkuId   string             `json:"skuId"`
	History []InventoryHistory `gorm:"ForeignKey:InventoryId" json:"history" form:"history"`

	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	CreatedAt time.Time  `json:"-"`
}

type InventoryFilter struct {
	SkuId string `form:"skuId"`
}

func (inf InventoryFilter) SetCriteria(db *gorm.DB) *gorm.DB {
	if len(inf.SkuId) == 0 {
		return db
	}
	return db.Where("sku_id like ?", SqlLikeString(inf.SkuId))
}
