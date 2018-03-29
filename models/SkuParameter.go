package models

import "time"

type SkuParameter struct {
	Id    uint   `gorm:"primary_key" json:"id" form:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
	SkuId string `json:"skuId"`

	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	BaseModel
}
