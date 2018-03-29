package middlewares

import (
	"github.com/aaa59891/go_demo_fullstack/constants"
	"github.com/aaa59891/go_demo_fullstack/controllers"
	"github.com/aaa59891/go_demo_fullstack/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	role, ok := session.Get(constants.SessionRole).(*models.Role)

	if !ok {
		c.Redirect(http.StatusFound, "/")
		return
	}

	if !checkAuth(*role, c.Request) {
		controllers.GoToErrorPage(http.StatusUnauthorized, c, constants.ErrNoAuth)
		return
	}

	c.Next()
}

func checkAuth(role models.Role, r *http.Request) bool {
	for _, rf := range role.Functions {
		function := rf.Function
		if function.Method == r.Method && function.Uri == r.URL.Path {
			return true
		}
	}
	return false
}
