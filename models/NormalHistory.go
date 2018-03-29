package models

import (
	"github.com/aaa59891/mosi_demo_go/utils"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	SaveAction     = "SAVE"
	DeleteAction   = "DELETE"
	ChangePassword = "CHANGE PASSWORD"
	SkuUploadFile  = "UPLOAD FILE"
	LoginAction    = "LOGIN"
)

var InventoryAction = map[int]string{
	AddLevel:   "STOCK",
	MinusLevel: "WITHDRAW",
}

var deleteMapKey = []string{
	"createdAt",
	"updatedAt",
	"deletedAt",
	"user",
}

type ModelJson map[string]interface{}

type NormalHistory struct {
	Id        uint      `gorm:"primary key"`
	ModelName string    `json:"modelName"`
	JsonData  ModelJson `sql:"TYPE:json" json:"jsonData"`
	Action    string    `json:"action"`
	User      string    `json:"user"`
	From      string    `json:"from"`
	CreatedAt time.Time `json:"createdAt"`
}

func (p ModelJson) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *ModelJson) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}
	if err := json.Unmarshal(source, &i); err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}

func (nh *NormalHistory) SetDataFromModel(model interface{}) error {
	if model == nil {
		return nil
	}

	nh.ModelName = utils.GetStructName(model)

	jsonMap, err := utils.ConvertToJsonMap(model)

	if err != nil {
		return err
	}

	utils.DeleteMapField(jsonMap, deleteMapKey...)

	nh.JsonData = jsonMap

	return nil
}

func (nh *NormalHistory) Insert(tx *gorm.DB) error {

	if err := tx.Create(nh).Error; err != nil {
		return err
	}

	return nil
}

func CreateNormalHistory(action, user, from string) NormalHistory {
	nh := NormalHistory{}
	nh.Action = action
	nh.User = user
	nh.From = from
	return nh
}
