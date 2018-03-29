package models

import (
	"github.com/aaa59891/mosi_demo_go/constants"
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	AddLevel   = 1
	MinusLevel = 0
)

type InventoryHistory struct {
	Id           uint       `gorm:"primary_key"`
	User         string     `json:"user"`
	Action       int        `json:"action"`
	Amount       int        `json:"amount"`
	InventoryId  uint       `json:"inventoryId"`
	CurrentLevel int        `json:"currentLevel"`
	DeletedAt    *time.Time `sql:"index" json:"deletedAt"`
	BaseModel
}

func (ir *InventoryHistory) InsertNewRecord(inventory *Inventory) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		var modifier int

		switch ir.Action {
		case AddLevel:
			modifier = 1
		case MinusLevel:
			modifier = -1
		default:
			return errors.New("Undefined action.")
		}

		updates := map[string]interface{}{
			"level": gorm.Expr("level + (? * ?) ", ir.Amount, modifier),
			"user":  ir.User,
		}

		if err := tx.Model(&Inventory{}).Where("id = ?", ir.InventoryId).Updates(updates).Error; err != nil {
			return err
		}

		if err := tx.Find(&inventory, "id = ? ", ir.InventoryId).Error; err != nil {
			return err
		}

		if inventory.Level < 0 {
			return errors.New("Stock is not enough.")
		}

		ir.CurrentLevel = inventory.Level

		if err := tx.Create(ir).Error; err != nil {
			return err
		}

		nh, err := CreateNormalHistoryByModel(ir, InventoryAction[ir.Action], ir.User, constants.FromBrowser)

		if err != nil {
			return err
		}

		if err := nh.Insert(tx); err != nil {
			return err
		}

		return nil
	}
}
