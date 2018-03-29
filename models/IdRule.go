package models

import (
	"fmt"
	"strconv"
	"time"
)

const (
	RuleCategory_Job = 1
)

type IdRule struct {
	Id          uint   `gorm:"primary key" json:"id" form:"id"`
	Name        string `json:"name"`
	Prefix      string `json:"prefix"`
	DigitNumber int    `json:"digitNumber"`
	IsDateShow  bool   `json:"isDateShow"`
	DateFormat  string
	Category    int `json:"category" form:"category"`

	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
}

func (ir IdRule) Example() string {
	now := time.Now()
	return ir.CreateId(now, 1)
}

func (ir IdRule) CreateId(t time.Time, index int) string {
	result := ir.Prefix
	if ir.IsDateShow {
		result += t.Format(ir.DateFormat)
	}
	result += fmt.Sprintf("%0"+strconv.Itoa(ir.DigitNumber)+"d", index)
	return result
}
