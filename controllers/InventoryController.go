package controllers

import (
	"github.com/aaa59891/mosi_demo_go/constants"
	"github.com/aaa59891/mosi_demo_go/db"
	"github.com/aaa59891/mosi_demo_go/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetInventoriesView(c *gin.Context) {
	current, err := GetCurrentPage(c)
	inf := models.InventoryFilter{}
	if err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}
	if err := c.BindQuery(&inf); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}
	inventories := []models.Inventory{}
	db := db.DB
	db = GetFilterCriteria(inf, db)
	pg, err := Pagination(10, current, db, &inventories)
	if err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
	}
	c.HTML(http.StatusOK, "inventoryList.html", gin.H{
		"data": inventories,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
		"pg":     pg,
		"filter": inf,
	})
}

func UpdateInventoryLevel(c *gin.Context) {
	inventoryRecord := models.InventoryHistory{}

	if err := c.BindJSON(&inventoryRecord); err != nil {
		GoToErrorPage(http.StatusBadRequest, c, err)
		return
	}

	inventory := models.Inventory{}

	inventoryRecord.User = GetSessionLoginUser(c)

	if err := models.Transactional(inventoryRecord.InsertNewRecord(&inventory)); err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	SendSocket(constants.SocketRoomMosi, constants.SocketEventInventoryUpdate, inventory)

	c.JSON(http.StatusOK, gin.H{
		"data": inventory,
	})
}

func GetInventoryHistoryView(c *gin.Context) {
	id := c.Query("id")
	inventory := models.Inventory{}

	if err := db.DB.Preload("History").Find(&inventory, "id = ? ", id).Error; err != nil {
		GoToErrorPage(http.StatusInternalServerError, c, err)
		return
	}

	c.HTML(http.StatusOK, "inventoryHistory.html", gin.H{
		"data": inventory,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
	})
}
