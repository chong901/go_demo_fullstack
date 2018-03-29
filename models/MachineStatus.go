package models

type MachineStatus struct {
	Id   int    `gorm:"primary_key" json:"id" form:"id"`
	Name string `json:"name"`

	BaseModel
}
