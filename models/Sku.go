package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Sku struct {
	Id             string  `gorm:"primary_key" json:"id" form:"id"`
	Ct             float64 `sql:"type:decimal(6,2);" json:"ct,string"`
	PerUnitMeasure float64 `sql:"type:decimal(8,2);" json:"perUnitMeasure,string"`
	Unit           string  `json:"unit"`
	GsUri          string
	SopUri         string
	User           string         `json:"user"`
	Parameters     []SkuParameter `gorm:"ForeignKey:SkuId;save_associations:false" json:"parameters" form:"parameters"`
	Inventory      Inventory      `form:"inventory"`

	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	BaseModel
}

type SkuFilter struct {
	Id string `form:"id"`
}

var ErrInventoryNotEmpty = errors.New("Can not delete this sku, please check the inventory.")

func (sku *Sku) Save(from string) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		inventory := Inventory{}

		if err := tx.Omit("created_at", "gs_uri", "sop_uri").Save(sku).Error; err != nil {
			return err
		}

		for _, parameter := range sku.Parameters {
			if err := tx.Omit("created_at").Save(&parameter).Error; err != nil {
				return err
			}
		}

		if err := tx.Find(&inventory, "sku_id = ?", sku.Id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				inventory.SkuId = sku.Id
				if err := tx.Create(&inventory).Error; err != nil {
					return err
				}

				nh, err := CreateNormalHistoryByModel(inventory, SaveAction, sku.User, from)

				if err != nil {
					return err
				}

				if err := nh.Insert(tx); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		nh, err := CreateNormalHistoryByModel(sku, SaveAction, sku.User, from)

		if err != nil {
			return err
		}

		if err := nh.Insert(tx); err != nil {
			return err
		}

		return nil
	}
}

func (sku *Sku) Remove(tx *gorm.DB) error {
	inventory := Inventory{}

	if err := tx.Find(&inventory, "sku_id = ?", sku.Id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	if inventory.Level > 0 {
		return ErrInventoryNotEmpty
	}

	if err := tx.Delete(sku).Delete(&inventory, "id = ?", inventory.Id).Error; err != nil {
		return err
	}

	return nil
}

func (sf SkuFilter) SetCriteria(db *gorm.DB) *gorm.DB {
	if len(sf.Id) == 0 {
		return db
	}
	return db.Where("id like ?", SqlLikeString(sf.Id))
}
