package models

import "time"

type RoleFunction struct {
	Id         uint     `gorm:"primary_key" json:"id" form:"id"`
	RoleId     int      `json:"roleId"`
	Function   Function `gorm:"ForeignKey:FunctionId;save_associations:false" form:"function"`
	FunctionId uint     `json:"functionId"`

	CreatedAt time.Time `json:"-"`
}
