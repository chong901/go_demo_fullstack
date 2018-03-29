package controllers

import (
	"github.com/aaa59891/mosi_demo_go/constants"
	"github.com/aaa59891/mosi_demo_go/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

func GetFunctionsSessionApi(c *gin.Context) {

	role, ok := GetSessionRole(c)

	if !ok {
		c.JSON(http.StatusOK, nil)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"roleId":    role.Id,
			"functions": getFunctions(role.Functions),
		})
	}

}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(constants.SessionLoginUser)
	session.Delete(constants.SessionRole)
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

func ChangeLang(c *gin.Context) {
	lang, _ := c.GetQuery("lang")
	session := sessions.Default(c)
	session.Set(constants.SessionCookieLang, lang)
	session.Save()
	c.JSON(http.StatusOK, nil)
}

func getFunctions(rfs []models.RoleFunction) []models.Function {
	functions := make(models.Functions, 0)
	for _, rf := range rfs {
		if rf.Function.IsMenu {
			functions = append(functions, rf.Function)
		}
	}
	sort.Sort(functions)
	return functions
}
