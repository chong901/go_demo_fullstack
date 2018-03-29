package models

import (
	"github.com/aaa59891/mosi_demo_go/constants"
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

var (
	ErrAccountNotEmpty = errors.New("Some accounts still belong to this role, cannot delete this role.")
)

type Role struct {
	Id        int            `gorm:"primary_key" json:"id" form:"id"`
	Name      string         `json:"name"`
	Functions []RoleFunction `gorm:"ForeignKey:RoleId;save_associations:false" json:"functions" form:"functions"`

	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	BaseModel
}

func (role *Role) SaveFunctions(tx *gorm.DB) error {

	if role.Id == 0 {
		return constants.ErrDataNotExist("role")
	}

	if err := tx.Delete(&RoleFunction{}, "role_id = ?", role.Id).Error; err != nil {
		return err
	}

	if len(role.Functions) > 0 {
		for _, function := range role.Functions {
			if err := tx.Create(&function).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (role *Role) Remove(tx *gorm.DB) error {
	if role.Id == 0 {
		return constants.ErrDataNotExist("role")
	}

	var count int

	if err := tx.Model(&Account{}).Where("role_id = ?", role.Id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return ErrAccountNotEmpty
	}

	if err := tx.Delete(&Role{}, "id = ?", role.Id).Error; err != nil {
		return err
	}

	if err := tx.Delete(&RoleFunction{}, "role_id = ?", role.Id).Error; err != nil {
		return err
	}

	return nil
}
