package models

import (
	"strings"
	"time"

	"github.com/aaa59891/go_demo_fullstack/db"
	"github.com/jinzhu/gorm"
)

const (
	JobStatusDone       = 1
	JobStatusUndone     = 2
	JobStatusPlanned    = 3
	JobStatusProcessing = 4
)

type Job struct {
	Id        string         `gorm:"primary_key" json:"id" form:"id"`
	SkuId     string         `json:"skuId"`
	Sku       Sku            `gorm:"ForeignKey:SkuId"`
	Quantity  int            `json:"quantity"`
	Lot       string         `json:"lot"`
	User      string         `json:"user"`
	MachineId uint           `json:"machineId"`
	Machine   Machine        `gorm:"ForeignKey:JobId" json:"machine" form:"machine"`
	IdRuleId  int            `json:"idRuleId"`
	Status    int            `json:"status"`
	Records   []RecordOutput `gorm:"ForeignKey:JobId"`

	NextJobId string     `json:"nextJobId"`
	JobOrder  int        `json:"jobOrder"`
	StartTime *time.Time `json:"startTime"`
	EndTime   *time.Time `json:"endTime"`

	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`

	OutputRecord []RecordOutput `gorm:"ForeignKey:JobId" json:"outputRecord" form:"outputRecord"`
}

type JobSchedule struct {
	UnplannedJobIds []string     `json:"unplannedJobIds"`
	PlannedJobs     []plannedJob `json:"plannedJobs"`
}

type plannedJob struct {
	Id        string `json:"id"`
	NextJobId string `json:"nextJobId"`
	MachineId int    `json:"machineId"`
	JobOrder  int    `json:"jobOrder"`
	Status    int    `json:"status"`
}

type JobFilter struct {
	Id        string    `form:"id"`
	SkuId     string    `form:"skuId"`
	Lot       string    `form:"lot"`
	StartDate time.Time `form:"startDate" time_format:"2006-01-02"`
	EndDate   time.Time `form:"endDate" time_format:"2006-01-02"`
}

func GetJobsByStatus(machineId, status int, orderString string) (jobs []Job, err error) {
	whereCondition := []interface{}{}
	if machineId != 0 {
		whereCondition = append(whereCondition, "machine_id = ? and status = ?")
		whereCondition = append(whereCondition, machineId)
	} else {
		whereCondition = append(whereCondition, "status = ?")
	}
	whereCondition = append(whereCondition, status)
	if len(orderString) > 0 {
		err = db.DB.Order(orderString).Find(&jobs, whereCondition...).Error
	} else {
		err = db.DB.Find(&jobs, whereCondition...).Error
	}
	return
}

func (js JobSchedule) SaveSchedule(tx *gorm.DB) error {
	updateUnplannedJob := map[string]interface{}{
		"machine_id":  0,
		"status":      JobStatusUndone,
		"next_job_id": "",
		"job_order":   0,
	}

	if err := tx.Model(&Job{}).Where("id in (?)", js.UnplannedJobIds).Update(updateUnplannedJob).Error; err != nil {
		return err
	}

	for _, pj := range js.PlannedJobs {
		status := JobStatusPlanned
		if pj.Status == JobStatusProcessing {
			status = JobStatusProcessing
		}
		updatePlannedJob := map[string]interface{}{
			"next_job_id": pj.NextJobId,
			"machine_id":  pj.MachineId,
			"status":      status,
			"job_order":   pj.JobOrder,
		}
		if err := tx.Model(&Job{}).Where("id = ?", pj.Id).Update(updatePlannedJob).Error; err != nil {
			return err
		}
	}

	return nil
}

func (jf JobFilter) SetCriteria(db *gorm.DB) *gorm.DB {
	whereStrSlice := make([]string, 0)
	whereValue := make([]interface{}, 0)
	if len(jf.Id) != 0 {
		whereStrSlice = append(whereStrSlice, "id like ?")
		whereValue = append(whereValue, SqlLikeString(jf.Id))
	}
	if len(jf.SkuId) > 0 {
		whereStrSlice = append(whereStrSlice, "sku_id like ?")
		whereValue = append(whereValue, SqlLikeString(jf.SkuId))
	}
	if len(jf.Lot) > 0 {
		whereStrSlice = append(whereStrSlice, "lot like ?")
		whereValue = append(whereValue, SqlLikeString(jf.Lot))
	}

	if !jf.EndDate.IsZero() || !jf.StartDate.IsZero() {
		jf.StartDate, jf.EndDate = GetFilterDateValue(jf.StartDate, jf.EndDate)
		whereStrSlice = append(whereStrSlice, "created_at between ? and ?")
		whereValue = append(whereValue, jf.StartDate, jf.EndDate)
	}

	return db.Where(strings.Join(whereStrSlice, " and "), whereValue...)
}

func (rf RecordSummaryFilter) SetCriteria(db *gorm.DB) *gorm.DB {
	whereStrSlice := make([]string, 0)
	whereValue := make([]interface{}, 0)
	if len(rf.Job) != 0 {
		whereStrSlice = append(whereStrSlice, "id like ?")
		whereValue = append(whereValue, SqlLikeString(rf.Job))
	}
	if len(rf.Sku) > 0 {
		whereStrSlice = append(whereStrSlice, "sku_id like ?")
		whereValue = append(whereValue, SqlLikeString(rf.Sku))
	}
	return db.Preload("Records").Where("status > ?", 2).Where(strings.Join(whereStrSlice, " and "), whereValue...)
}
