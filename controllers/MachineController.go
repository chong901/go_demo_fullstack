package controllers

import (
	"net/http"

	"github.com/aaa59891/go_demo_fullstack/constants"
	"github.com/aaa59891/go_demo_fullstack/db"
	"github.com/aaa59891/go_demo_fullstack/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type MachineInfo struct {
	Name   string
	Status int
	Job    string
	Output int
	Qty    int
	Sku    string
	Ct     float64
}

func MachineListView(c *gin.Context) {
	currentPage, err := GetCurrentPage(c)
	if err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	mf := models.MachineFilter{}
	if err := c.BindQuery(&mf); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	machines := []models.Machine{}
	db := GetFilterCriteria(mf, db.DB)
	pg, err := Pagination(10, currentPage, db, &machines)

	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.HTML(http.StatusOK, "machineM.html", gin.H{
		"data": machines,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
		"pg":     pg,
		"filter": mf,
	})
}

func EditMachineView(c *gin.Context) {
	id := c.Query("id")
	machine := models.Machine{}

	if len(id) != 0 {
		if err := db.DB.Find(&machine, "id = ?", id).Error; err != nil {
			GoToErrorPage(http.StatusInternalServerError, c, err)
			return
		}
	}

	c.HTML(http.StatusOK, "machineMEdit.html", gin.H{
		"data": machine,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
	})
}

func CreateMachine(c *gin.Context) {
	machine := models.Machine{}

	if err := c.BindJSON(&machine); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	machine.User = GetSessionLoginUser(c)

	if err := models.Transactional(
		func(tx *gorm.DB) error {
			if err := machine.InsertNewData(tx); err != nil {
				return err
			}

			nh, err := models.CreateNormalHistoryByModel(machine, models.SaveAction, GetSessionLoginUser(c), constants.FromBrowser)

			if err != nil {
				return err
			}

			return nh.Insert(tx)
		}); err != nil {
		if err == constants.ErrIdDuplicated {
			GoToErrorPage(http.StatusBadRequest, c, err)
		} else {
			GoToErrorPage(http.StatusInternalServerError, c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": machine.Id,
	})
}

func UpdateMachine(updateColumns []interface{}, from string) gin.HandlerFunc {
	return func(c *gin.Context) {
		machine := models.Machine{}

		if err := c.BindJSON(&machine); err != nil {
			GoToErrorPage(http.StatusBadRequest, c, err)
			return
		}

		machine.User = GetSessionLoginUser(c)

		switch from {
		case constants.FromBrowser:

			if err := models.Transactional(
				func(tx *gorm.DB) error {
					if err := machine.Update(updateColumns, from)(tx); err != nil {
						return err
					}

					nh, err := models.CreateNormalHistoryByModel(machine, models.SaveAction, GetSessionLoginUser(c), constants.FromBrowser)

					if err != nil {
						return err
					}

					return nh.Insert(tx)
				}); err != nil {
				GoToErrorPage(http.StatusInternalServerError, c, err)
				return
			}
		case constants.FromPi:
			if err := models.Transactional(machine.Update(updateColumns, from)); err != nil {
				GoToErrorPage(http.StatusInternalServerError, c, err)
				return
			}
		}

		SendSocket("mosi", "machineUpdate", machine)

		c.JSON(http.StatusOK, gin.H{
			"id": machine.Id,
		})
	}
}

func DeleteMachine(c *gin.Context) {
	machine := models.Machine{}

	if err := c.BindQuery(&machine); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	nh, err := models.CreateNormalHistoryByModel(machine, models.DeleteAction, GetSessionLoginUser(c), constants.FromBrowser)

	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	if err := models.Transactional(
		nh.Insert,
		models.DeleteById(&models.Machine{}, machine.Id),
	); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func GetMachineSummaryApi(c *gin.Context) {
	machineSummary := []MachineInfo{}
	err := db.DB.Raw(
		`SELECT machines.name, machines.machine_status_id as status, machines.job_id as job, machines.output_count_total as output, jobs.quantity as qty, jobs.sku_id as sku, skus.ct as ct
		FROM (select * from machines where deleted_at is null) as machines
		LEFT JOIN jobs ON machines.job_id = jobs.id LEFT JOIN skus on skus.id = jobs.sku_id
		ORDER BY machines.id`).Scan(&machineSummary).Error
	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": machineSummary,
	})
}
