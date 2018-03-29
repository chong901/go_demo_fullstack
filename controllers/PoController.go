package controllers

import (
	"github.com/aaa59891/mosi_demo_go/constants"
	"github.com/aaa59891/mosi_demo_go/db"
	"github.com/aaa59891/mosi_demo_go/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func EditPoView(c *gin.Context) {
	id := c.Query("id")

	po := models.Po{}
	skus := []models.Sku{}
	units := []models.Configuration{}

	if len(id) != 0 {
		if err := db.DB.Preload("SkuList").Find(&po, "id = ?", id).Error; err != nil {
			GoToErrorPage(http.StatusInternalServerError, c, err)
			return
		}
	}

	if err := db.DB.Find(&skus).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	if err := db.DB.Find(&units, "type = ?", models.Unit).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.HTML(http.StatusOK, "poMEdit.html", gin.H{
		"data":                    po,
		"skus":                    skus,
		"units":                   units,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
	})
}

func SavePo(c *gin.Context) {
	po := models.Po{}

	if err := c.BindJSON(&po); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}
	user := GetSessionLoginUser(c)

	po.User = user

	if err := models.Transactional(
		func(tx *gorm.DB) error {
			if err := po.Save(tx); err != nil {
				return err
			}

			nh, err := models.CreateNormalHistoryByModel(po, models.SaveAction, user, constants.FromBrowser)

			if err != nil {
				return err
			}

			return nh.Insert(tx)
		}); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	SendSocket(constants.SocketRoomMosi, constants.SocketEventPoUpdate, po)

	c.JSON(http.StatusOK, gin.H{
		"id": po.Id,
	})
}

func GetPosView(c *gin.Context) {
	pos := []models.Po{}

	if err := db.DB.Preload("SkuList").Find(&pos).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.HTML(http.StatusOK, "poMList.html", gin.H{
		"data": pos,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
	})
}

func DeletePo(c *gin.Context) {
	po := models.Po{}
	if err := c.BindQuery(&po); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	nh, err := models.CreateNormalHistoryByModel(po, models.DeleteAction, GetSessionLoginUser(c), constants.FromBrowser)

	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	if err := models.Transactional(
		nh.Insert,
		models.DeleteById(&models.Po{}, po.Id),
	); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func DeletePoSku(c *gin.Context) {
	ps := models.PoSku{}

	if err := c.BindQuery(&ps); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	if err := db.DB.Find(&ps).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	nh, err := models.CreateNormalHistoryByModel(ps, models.DeleteAction, GetSessionLoginUser(c), constants.FromBrowser)

	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	if err := models.Transactional(
		nh.Insert,
		models.DeleteById(&models.PoSku{}, ps.Id),
	); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
