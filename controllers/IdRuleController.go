package controllers

import (
	"github.com/aaa59891/mosi_demo_go/constants"
	"github.com/aaa59891/mosi_demo_go/db"
	"github.com/aaa59891/mosi_demo_go/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var backUrlMap = map[int]string{
	models.RuleCategory_Job: "/jobIdRuleList",
}

var addUrlMap = map[int]string{
	models.RuleCategory_Job: "/jobIdRule",
}

func IdRuleListView(category int) func(c *gin.Context) {
	return func(c *gin.Context) {
		models := []models.IdRule{}

		if err := db.DB.Find(&models, "category = ?", category).Error; err != nil {
			GoToErrorPage(http.StatusInternalServerError, c, err)
			return
		}

		c.HTML(http.StatusOK, "IdRuleList.html", gin.H{
			"data":                    models,
			"editUrl":                 addUrlMap[category],
			constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
		})
	}
}

/**
Models' id needs to add tag `form:"id"`
*/
func IdRuleView(category int) func(c *gin.Context) {
	return func(c *gin.Context) {
		model := models.IdRule{}

		if err := c.BindQuery(&model); err != nil {
			GoToErrorPage(http.StatusBadRequest, c, err)
			return
		}

		if !db.DB.NewRecord(model) {
			if err := db.DB.Find(&model, "category = ?", category).Error; err != nil {
				GoToErrorPage(http.StatusInternalServerError, c, err)
				return
			}
		} else {
			model.Category = category
		}
		fmt.Println(model)

		c.HTML(http.StatusOK, "IdRuleEdit.html", gin.H{
			"data":                    model,
			"backUrl":                 backUrlMap[category],
			constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
		})
	}
}

func CreateIdRule(c *gin.Context) {
	model := models.IdRule{}

	if err := c.BindJSON(&model); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if model.IsDateShow {
		model.DateFormat = "20060102"
	}

	if err := db.DB.Create(&model).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": model,
	})
}

func UpdateIdRule(c *gin.Context) {
	model := models.IdRule{}

	if err := c.BindJSON(&model); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if model.IsDateShow {
		model.DateFormat = "20060102"
	}

	if err := db.DB.Select(nil, "name", "is_date_show", "prefix", "digit_number", "date_format").Save(&model).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": model,
	})
}

func DeleteIdRule(category int) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Query("id")

		if err := db.DB.Delete(&models.IdRule{}, "id = ? and category = ?", id, category).Error; err != nil {
			GoToErrorPage(http.StatusInternalServerError, c, err)
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func IdRuleListApi(category int) func(c *gin.Context) {
	return func(c *gin.Context) {
		ruleList := []models.IdRule{}

		if err := db.DB.Find(&ruleList, "category = ?", category).Error; err != nil {
			GoToErrorPage(http.StatusInternalServerError, c, err)
			return
		}

		type tempData struct {
			Id      uint   `json:"id"`
			Name    string `json:"name"`
			Example string `json:"example"`
		}

		dataList := make([]tempData, 0)

		for _, rule := range ruleList {
			temp := tempData{rule.Id, rule.Name, rule.Example()}
			dataList = append(dataList, temp)
		}

		c.JSON(http.StatusOK, gin.H{
			"data": dataList,
		})
	}
}
