package middlewares

import (
	"github.com/aaa59891/mosi_demo_go/constants"
	"github.com/aaa59891/mosi_demo_go/controllers"
	"github.com/gin-gonic/gin"
)

func I18nTranslate(c *gin.Context) {
	cookieLang := controllers.GetSessionString(c, constants.SessionCookieLang)
	headerLang := controllers.GetSessionString(c, constants.SessionHeaderLang)
	c.Set(constants.ContextSetLang, cookieLang)
	if len(cookieLang) == 0 {
		c.Set(constants.ContextSetLang, headerLang)
	}
	c.Next()
}
