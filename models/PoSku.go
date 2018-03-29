package models

import "time"

type PoSku struct {
	Id       uint   `gorm:"primary_key" json:"id" form:"id"`
	SkuId    string `json:"skuId"`
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
	PoId     string `json:"poId"`

	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	BaseModel
}
