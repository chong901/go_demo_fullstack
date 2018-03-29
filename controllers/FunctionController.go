package controllers

import (
	"net/http"

	"github.com/aaa59891/go_demo_fullstack/constants"
	"github.com/aaa59891/go_demo_fullstack/db"
	"github.com/aaa59891/go_demo_fullstack/models"
	"github.com/gin-gonic/gin"
)

func GetFuncsView(c *gin.Context) {
	funcs, err := models.GetFunctionAll()

	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.HTML(http.StatusOK, "functionList.html", gin.H{
		"data": funcs,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
	})
}

func GetFuncsDictApi(c *gin.Context) {
	fd := []models.FuncsDict{}

	if err := db.DB.Table("functions").Select("id, name").Where("deleted_at is null").Scan(&fd).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": fd,
	})
}

func GetFuncByIdApi(c *gin.Context) {
	fun := models.Function{}

	if err := c.BindQuery(&fun); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := db.DB.First(&fun).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": fun,
	})
}

func SaveFunc(c *gin.Context) {
	fun := models.Function{}

	if err := c.BindJSON(&fun); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := fun.Save(); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": fun,
	})
}

func DeleteFunc(c *gin.Context) {
	function := models.Function{}

	if err := c.Bind(&function); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := models.Transactional(function.DeleteFunc); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
