package controllers

import (
	"net/http"
	"strconv"

	"github.com/aaa59891/go_demo_fullstack/constants"
	"github.com/aaa59891/go_demo_fullstack/db"
	"github.com/aaa59891/go_demo_fullstack/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type RecordOutputSummary struct {
	Job      string
	Sku      string
	Planned  int
	Quantity int
}

func GetRecordsView(c *gin.Context) {
	currentPage, err := GetCurrentPage(c)
	if err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}
	rf := models.RecordSummaryFilter{}
	if err := c.BindQuery(&rf); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	jobs := []models.Job{}
	db := GetFilterCriteria(rf, db.DB)
	pg, err := Pagination(10, currentPage, db, &jobs)

	summary := make([]RecordOutputSummary, len(jobs))
	for i, job := range jobs {
		summary[i].Job = job.Id
		summary[i].Sku = job.SkuId
		summary[i].Planned = job.Quantity
		for _, item := range job.Records {
			summary[i].Quantity += item.Quantity
		}
	}

	// if err := db.DB.Preload("Sku").Preload("Records").Where("status > ?", 2).Find(&jobs); err != nil {
	// 	GoToErrorPage(http.StatusInternalServerError, c, err)
	// 	return
	// }
	c.HTML(http.StatusOK, "recordOutput.html", gin.H{
		"data": summary,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
		"pg":     pg,
		"filter": rf,
	})
}

func GetRecordHistoryView(c *gin.Context) {
	records := []models.RecordOutput{}
	id := c.Query("id")
	sum := 0

	if len(id) != 0 {
		if err := db.DB.Find(&records, "job_id = ?", id).Error; err != nil {
			GoToErrorPage(http.StatusInternalServerError, c, err)
			return
		}
		for _, item := range records {
			sum += item.Quantity
		}
	}

	c.HTML(http.StatusOK, "recordOutputHistory.html", gin.H{
		"job":      id,
		"totalQty": sum,
		"data":     records,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
	})
}

func CreateRecord(c *gin.Context) {
	record := models.RecordOutput{}

	if err := c.BindJSON(&record); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	record.User = GetSessionLoginUser(c)

	if err := db.DB.Create(&record).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	SendSocket(constants.SocketRoomMosi, constants.SocketEventReportOutputUpdate, record)
	c.JSON(http.StatusOK, gin.H{
		"data": record,
	})
}

func UpdateRecord(c *gin.Context) {
	record := models.RecordOutput{}
	if err := c.BindJSON(&record); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	val, err := strconv.ParseUint(c.Query("id"), 10, 0)
	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}
	record.Id = uint(val)
	record.User = GetSessionLoginUser(c)

	if err := models.Transactional(
		func(tx *gorm.DB) error {
			if err := tx.Model(record).Updates(map[string]interface{}{"quantity": record.Quantity, "user": record.User}).Error; err != nil {
				return err
			}
			nh, err := models.CreateNormalHistoryByModel(record, models.SaveAction, GetSessionLoginUser(c), constants.FromBrowser)
			if err != nil {
				return err
			}
			return nh.Insert(tx)
		}); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	SendSocket(constants.SocketRoomMosi, constants.SocketEventReportOutputUpdate, record)
	c.JSON(http.StatusOK, gin.H{
		"data": record,
	})
}

func DeleteRecord(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, "ID not specified.")
		return
	}

	if err := db.DB.Delete(&models.RecordOutput{}, "id = ?", id).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}
