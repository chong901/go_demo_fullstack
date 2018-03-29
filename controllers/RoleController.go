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

func GetRoleSettingView(c *gin.Context) {
	funcs, err := models.GetFunctionAll()

	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.HTML(http.StatusOK, "roleSetting.html", gin.H{
		"data": funcs,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
	})
}

func GetRolesApi(c *gin.Context) {
	roles := []models.Role{}

	role, _ := GetSessionRole(c)

	if err := db.DB.Find(&roles, "id >= ?", role.Id).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": roles,
	})
}

func GetRoleApi(c *gin.Context) {
	role := models.Role{}
	id := c.Query("id")

	idNum, err := strconv.Atoi(id)

	if err != nil {
		GoToErrorPage(http.StatusBadRequest, c, constants.ErrParameters)
		return
	}

	roleS, _ := GetSessionRole(c)

	if idNum < roleS.Id {
		GoToErrorPage(http.StatusUnauthorized, c, constants.ErrNoAuth)
		return
	}

	if err := db.DB.Preload("Functions").Find(&role, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			GoToErrorPage(http.StatusNotFound, c, err)
		} else {
			GoToErrorPage(http.StatusInternalServerError, c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": role,
	})
}

func CreateNewRole(c *gin.Context) {
	role := models.Role{}

	if err := c.BindJSON(&role); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := models.Transactional(
		func(tx *gorm.DB) error {
			if err := tx.Create(&role).Error; err != nil {
				return err
			}

			nh, err := models.CreateNormalHistoryByModel(role, models.SaveAction, GetSessionLoginUser(c), constants.FromBrowser)

			if err != nil {
				return err
			}

			return nh.Insert(tx)
		},
	); err != nil {
		if constants.IsDuplicatedErr(err) {
			GoToErrorPage(http.StatusBadRequest, c, constants.ErrIdDuplicated)
		} else {
			GoToErrorPage(http.StatusInternalServerError, c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": role,
	})
}

func UpdateRoleName(c *gin.Context) {
	role := models.Role{}

	if err := c.BindJSON(&role); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := models.Transactional(
		func(tx *gorm.DB) error {
			if err := tx.Model(&role).Update("name", role.Name).Error; err != nil {
				return err
			}

			nh, err := models.CreateNormalHistoryByModel(role, models.SaveAction, GetSessionLoginUser(c), constants.FromBrowser)

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
		"data": role,
	})
}

func UpdateRoleFunctions(c *gin.Context) {
	role := models.Role{}

	if err := c.BindJSON(&role); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := models.Transactional(
		func(tx *gorm.DB) error {
			if err := role.SaveFunctions(tx); err != nil {
				return err
			}

			nh, err := models.CreateNormalHistoryByModel(role, models.SaveAction, GetSessionLoginUser(c), constants.FromBrowser)

			if err != nil {
				return err
			}

			return nh.Insert(tx)
		}); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func DeleteRole(c *gin.Context) {
	role := models.Role{}

	if err := c.BindQuery(&role); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := models.Transactional(role.Remove); err != nil {
		if err == models.ErrAccountNotEmpty {
			GoToErrorPage(http.StatusBadRequest, c, err)
		} else {
			GoToErrorPage(http.StatusInternalServerError, c, err)
		}
		return
	}

	c.JSON(http.StatusOK, nil)
}
