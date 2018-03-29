package controllers

import (
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/aaa59891/mosi_demo_go/constants"
	"github.com/aaa59891/mosi_demo_go/controllers/filter"
	"github.com/aaa59891/mosi_demo_go/controllers/pagination"
	"github.com/aaa59891/mosi_demo_go/inits"
	"github.com/aaa59891/mosi_demo_go/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
	"github.com/jinzhu/gorm"
)

func GoToErrorPage(statusCode int, c *gin.Context, err error) {
	defer log.Panic(err)
	ajax := c.GetHeader("X-Requested-With")
	errMsg := err.Error()

	if statusCode == http.StatusInternalServerError {
		errMsg = constants.ErrMsgInternalServerError
	}
	if len(ajax) != 0 {
		c.JSON(statusCode, gin.H{
			"message": errMsg,
		})
		return
	}
	c.HTML(statusCode, "errorPage.html", gin.H{
		"message":                 errMsg,
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
	})
}

func GetSessionRole(c *gin.Context) (models.Role, bool) {
	session := sessions.Default(c)

	role, ok := session.Get(constants.SessionRole).(*models.Role)
	if ok {
		return *role, ok
	} else {
		return models.Role{}, ok
	}
}

func GetSessionLoginUser(c *gin.Context) string {
	loginUser, _ := getSessionValue(c, constants.SessionLoginUser).(string)
	return loginUser
}

func GetSessionString(c *gin.Context, key string) string {
	str, ok := getSessionValue(c, key).(string)
	if !ok {
		return ""
	}
	return str
}

func getSessionValue(c *gin.Context, key string) interface{} {
	session := sessions.Default(c)

	return session.Get(key)
}

func SendSocket(room, event string, data interface{}) {
	dataJson, err := json.Marshal(data)
	if err != nil {
		log.Println("Parse json had an error, ", err)
	} else {
		inits.SocketServer.BroadcastTo(room, event, string(dataJson))
	}
}

func Pagination(numPerPage, currentPage int, prepareDb *gorm.DB, model interface{}) (pagination.Pagination, error) {
	var totalCount int
	pg := pagination.Pagination{}
	if err := prepareDb.Model(model).Count(&totalCount).Error; err != nil {
		return pg, err
	}

	totalPage := int(math.Ceil(float64(totalCount) / float64(numPerPage)))

	if err := prepareDb.Limit(numPerPage).Offset(numPerPage * (currentPage - 1)).Find(model).Error; err != nil {
		return pg, err
	}

	pg.Current = currentPage
	pg.Next = currentPage + 1
	pg.Previous = currentPage - 1
	pg.Total = totalPage
	return pg, nil
}

func GetFilterCriteria(filter filter.Filter, db *gorm.DB) *gorm.DB {
	return filter.SetCriteria(db)
}

func GetCurrentPage(c *gin.Context) (current int, err error) {
	currentStr := c.Query("page")
	if len(currentStr) == 0 {
		current = 1
		return
	}
	current, err = strconv.Atoi(currentStr)
	return
}
