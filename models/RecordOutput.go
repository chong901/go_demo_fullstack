package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type RecordOutput struct {
	Id        uint       `gorm:"primary_key" json:"id"`
	JobId     string     `json:"jobId"`
	Quantity  int        `json:"quantity"`
	User      string     `json:"user"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	BaseModel
}

type RecordSummaryFilter struct {
	Job string `form:"job"`
	Sku string `form:"sku"`
}

func (r *RecordOutput) Update(updateColumns map[string]interface{}) func(db2 *gorm.DB) error {
	return func(tx *gorm.DB) error {

		if err := tx.Model(r).Updates(updateColumns).Error; err != nil {
			return err
		}

		return nil
	}
}
