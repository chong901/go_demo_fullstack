package models

import (
	"github.com/aaa59891/go_demo_fullstack/constants"
	"github.com/aaa59891/go_demo_fullstack/db"
	"github.com/jinzhu/gorm"
	"time"
)

type Function struct {
	Id         uint      `gorm:"primary_key" json:"id" form:"id"`
	Parent     int       `json:"parent"`
	Uri        string    `json:"uri"`
	Method     string    `json:"method"`
	Name       string    `json:"name"`
	OrderNum   int       `json:"orderNum"`
	IsMenu     bool      `json:"isMenu"`
	ParentFunc *Function `gorm:"ForeignKey:Id;AssociationForeignKey:Parent;save_associations:false" json:"-" form:"parentFunc"`

	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	BaseModel
}

type FuncsDict struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (f Function) DeleteFunc(tx *gorm.DB) error {
	if f.Id == 0 {
		return constants.ErrDataNotExist("function")
	}

	if err := tx.Delete(&RoleFunction{}, "function_id = ?", f.Id).Error; err != nil {
		return err
	}

	if err := tx.Delete(&f, "id = ?", f.Id).Error; err != nil {
		return err
	}

	return nil
}

// Weird: this model uses db.DB.Omit("created_at").Save(f) directly will not insert created_at when inserting data
func (f *Function) Save() error {
	if f.Id == 0 {
		return db.DB.Create(f).Error
	}

	return db.DB.Omit("created_at").Save(f).Error
}

func GetFunctionAll() ([]Function, error) {
	funcs := []Function{}
	err := db.DB.Preload("ParentFunc").Order("parent, order_num").Find(&funcs).Error

	return funcs, err
}

type Functions []Function

func (slice Functions) Len() int {
	return len(slice)
}

func (slice Functions) Less(i, j int) bool {
	if slice[i].Parent == slice[j].Parent {
		return slice[i].OrderNum < slice[j].OrderNum
	}
	return slice[i].Parent < slice[j].Parent
}

func (slice Functions) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
