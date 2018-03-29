package models

import (
	"github.com/aaa59891/go_demo_fullstack/constants"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

var (
	UpdateColumnsFromBrowser = []interface{}{
		"name",
		"type",
		"department",
		"user",
		"ip_address",
	}

	UpdateColumnsFromPi = []interface{}{
		"machine_status_id",
		"worker_id",
		"job_id",
	}
)

type Machine struct {
	Id               uint          `gorm:"primary_key" json:"id" form:"id"`
	Name             string        `json:"name"`
	Type             string        `json:"type"`
	Department       string        `json:"department"`
	MachineStatusId  int           `json:"machineStatusId"`
	WorkerId         string        `json:"workerId"`
	JobId            string        `json:"jobId"`
	MachineStatus    MachineStatus `json:"machineStatus" form:"machineStatus"`
	User             string        `json:"user"`
	OutputCountTotal int           `json:"outputCountTotal"`
	IpAddress        string        `json:"ipAddress"`

	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"-"`
}

type MachineFilter struct {
	Name       string `form:"name"`
	Type       string `form:"type"`
	Department string `form:"department"`
}

func (m *Machine) InsertNewData(tx *gorm.DB) error {
	if !tx.NewRecord(m) {
		return constants.ErrIdDuplicated
	}

	if err := tx.Create(m).Error; err != nil {
		return err
	}

	return nil
}

func (m Machine) IsWorkerJobStatusExist() bool {
	return len(m.WorkerId) > 0 || len(m.JobId) > 0 || int(m.MachineStatusId) > 0
}

func (m *Machine) Update(updateColumns []interface{}, from string) func(db2 *gorm.DB) error {
	return func(tx *gorm.DB) error {

		if err := tx.Select(nil, updateColumns...).Save(m).Error; err != nil {
			return err
		}

		return nil
	}
}

func (mf MachineFilter) SetCriteria(db *gorm.DB) *gorm.DB {
	whereStrSlice := make([]string, 0)
	whereValue := make([]interface{}, 0)
	if len(mf.Name) > 0 {
		whereStrSlice = append(whereStrSlice, "name like ?")
		whereValue = append(whereValue, SqlLikeString(mf.Name))
	}

	if len(mf.Type) > 0 {
		whereStrSlice = append(whereStrSlice, "type like ?")
		whereValue = append(whereValue, SqlLikeString(mf.Type))
	}

	if len(mf.Department) > 0 {
		whereStrSlice = append(whereStrSlice, "department like ?")
		whereValue = append(whereValue, SqlLikeString(mf.Department))
	}

	return db.Where(strings.Join(whereStrSlice, " and "), whereValue...)
}
