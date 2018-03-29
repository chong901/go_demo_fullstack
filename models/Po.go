package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Po struct {
	Id       string    `gorm:"primary_key" json:"id" form:"id"`
	ClientId string    `json:"clientId"`
	Due      time.Time `json:"due"`
	Status   string    `json:"status"`
	SkuList  []PoSku   `gorm:"ForeignKey:PoId;save_associations:false" json:"skuList" form:"skuList"`

	User      string     `json:"user"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"-"`
}

func (po *Po) Save(tx *gorm.DB) error {

	if err := tx.Omit("created_at").Save(po).Error; err != nil {
		return err
	}

	for _, skuData := range po.SkuList {
		if err := tx.Omit("created_at").Save(&skuData).Error; err != nil {
			return err
		}
	}

	return nil
}
