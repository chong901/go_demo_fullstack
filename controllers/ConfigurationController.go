package controllers

import (
	"net/http"

	"github.com/aaa59891/go_demo_fullstack/constants"
	"github.com/aaa59891/go_demo_fullstack/db"
	"github.com/aaa59891/go_demo_fullstack/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func GetConfigurationByType(t int) gin.HandlerFunc {
	return func(c *gin.Context) {
		configurations := []models.Configuration{}

		if err := db.DB.Find(&configurations, "type = ?", t).Error; err != nil {
			GoToErrorPage(http.StatusInternalServerError, c, err)
			return
		}

		c.HTML(http.StatusOK, "configurations.html", gin.H{
			"title": models.ConfigurationTitleDict[t],
			"data":  configurations,
			"type":  t,
			constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
		})
	}
}

func GetConfigurationByTypeApi(t int) gin.HandlerFunc {
	return func(c *gin.Context) {
		configurations := []models.Configuration{}
		if err := db.DB.Find(&configurations, "type = ?", t).Error; err != nil {
			GoToErrorPage(http.StatusInternalServerError, c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": configurations,
		})
	}
}

func AddConfiguration(c *gin.Context) {
	cf := models.Configuration{}

	if err := c.Bind(&cf); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := models.Transactional(
		func(tx *gorm.DB) error {
			if err := tx.Create(&cf).Error; err != nil {
				return err
			}

			nh, err := models.CreateNormalHistoryByModel(cf, models.SaveAction, GetSessionLoginUser(c), constants.FromBrowser)

			if err != nil {
				return err
			}

			return nh.Insert(tx)
		},
	); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   cf.Id,
		"name": cf.Name,
	})
}

func DeleteConfiguration(c *gin.Context) {
	con := models.Configuration{}

	if err := c.BindQuery(&con); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	nh, err := models.CreateNormalHistoryByModel(con, models.DeleteAction, GetSessionLoginUser(c), constants.FromBrowser)

	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	if err := models.Transactional(
		nh.Insert,
		models.DeleteById(&models.Configuration{}, con.Id),
	); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
