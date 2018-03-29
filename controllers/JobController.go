package controllers

import (
	"github.com/aaa59891/mosi_demo_go/constants"
	"github.com/aaa59891/mosi_demo_go/db"
	"github.com/aaa59891/mosi_demo_go/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

type PartInfo struct {
	Id      string
	Planned int
	Output  int
}

func EditJobView(c *gin.Context) {
	id := c.Query("id")
	job := models.Job{}
	skus := []models.Sku{}
	idRuleList := []models.IdRule{}
	if len(id) != 0 {
		if err := db.DB.Find(&job, "id = ?", id).Error; err != nil {
			GoToErrorPage(http.StatusInternalServerError, c, err)
			return
		}
	} else {
		if err := db.DB.Find(&idRuleList, "category = ?", models.RuleCategory_Job).Error; err != nil {
			GoToErrorPage(http.StatusInternalServerError, c, err)
			return
		}
	}
	if err := db.DB.Find(&skus).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.HTML(http.StatusOK, "jobMEdit.html", gin.H{
		"data":                    job,
		"skus":                    skus,
		"idRuleList":              idRuleList,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
	})
}

func JobListView(c *gin.Context) {
	currentPage, err := GetCurrentPage(c)
	if err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	jf := models.JobFilter{}
	if err := c.BindQuery(&jf); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	jobs := []models.Job{}

	db := db.DB.Preload("Machine")
	db = GetFilterCriteria(jf, db)
	pg, err := Pagination(10, currentPage, db, &jobs)

	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.HTML(http.StatusOK, "jobMList.html", gin.H{
		"data": jobs,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
		"pg":     pg,
		"filter": jf,
	})
}

func SaveJob(c *gin.Context) {
	job := models.Job{}
	if err := c.BindJSON(&job); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	isNew := false
	if db.DB.NewRecord(job) && job.IdRuleId != 0 {
		idRule := models.IdRule{}
		var count int
		if err := db.DB.Find(&idRule, "id = ?", job.IdRuleId).Error; err != nil {
			GoToErrorPage(http.StatusInternalServerError, c, err)
			return
		}
		if err := db.DB.Model(&job).Where("id_rule_id = ?", job.IdRuleId).Count(&count).Error; err != nil {
			GoToErrorPage(http.StatusInternalServerError, c, err)
			return
		}
		count++
		job.Id = idRule.CreateId(time.Now(), count)
		job.Status = models.JobStatusUndone
		isNew = true
	}

	user := GetSessionLoginUser(c)

	job.User = user

	if err := models.Transactional(
		func(tx *gorm.DB) error {
			var err error

			if isNew {
				err = tx.Create(&job).Error
			} else {
				err = tx.Model(&job).Updates(map[string]interface{}{"sku_id": job.SkuId, "quantity": job.Quantity, "lot": job.Lot}).Error
			}

			if err != nil {
				return err
			}

			nh, err := models.CreateNormalHistoryByModel(job, models.SaveAction, user, constants.FromBrowser)

			if err != nil {
				return err
			}

			return nh.Insert(tx)
		},
	); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	SendSocket(constants.SocketRoomMosi, constants.SocketEventJobUpdate, job)

	c.JSON(http.StatusOK, gin.H{
		"data": job,
	})
}

func DeleteJob(c *gin.Context) {
	job := models.Job{}

	if err := c.BindQuery(&job); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	nh, err := models.CreateNormalHistoryByModel(job, models.DeleteAction, GetSessionLoginUser(c), constants.FromBrowser)

	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	if err := models.Transactional(
		nh.Insert,
		models.DeleteById(&models.Job{}, job.Id),
	); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func GetPartQtyApi(c *gin.Context) {
	partSummary := []PartInfo{}
	err := db.DB.Raw(
		"SELECT jobs.sku_id as id, sum(jobs.quantity) as planned, sum(machines.output_count_total) as output " +
			"FROM jobs left join machines on jobs.id = machines.job_id " +
			"WHERE jobs.end_time is NULL " +
			"GROUP BY jobs.sku_id").Scan(&partSummary).Error
	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, partSummary)
}

func JobScheduleView(c *gin.Context) {
	machines := []models.Machine{}
	if err := db.DB.Find(&machines).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.HTML(http.StatusOK, "arrangeJob.html", gin.H{
		"machines":                machines,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
	})
}

func SaveJobSchedule(c *gin.Context) {
	jobSchedule := models.JobSchedule{}

	if err := c.BindJSON(&jobSchedule); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := models.Transactional(jobSchedule.SaveSchedule); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func GetPlannedJobsApi(c *gin.Context) {
	machineId := c.Query("machineId")
	machineIdNum, err := strconv.Atoi(machineId)

	if err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	plannedJobs, err := models.GetJobsByStatus(machineIdNum, models.JobStatusPlanned, "job_order")

	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	processJobs, err := models.GetJobsByStatus(machineIdNum, models.JobStatusProcessing, "")
	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"plannedJobs": plannedJobs,
		"processJobs": processJobs,
		"totalLength": len(plannedJobs) + len(processJobs),
	})
}

func GetUndoneJobsApi(c *gin.Context) {
	undoneJobs, err := models.GetJobsByStatus(0, models.JobStatusUndone, "created_at")
	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"undoneJobs": undoneJobs,
	})
}
