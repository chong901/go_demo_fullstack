package models

const (
	Parameters = 1
	Unit       = 2
)

var ConfigurationTitleDict = map[int]string{
	Parameters: "tSkuParameter",
	Unit:       "tUnitSetting",
}

type Configuration struct {
	Id   uint   `gorm:"primary_key" form:"id" json:"id"`
	Name string `form:"name"`
	Type int    `form:"type"`

	BaseModel
}
