package controllers

import (
	"net/http"

	"github.com/aaa59891/go_demo_fullstack/constants"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	sessionCookieLang := GetSessionString(c, constants.SessionCookieLang)
	if len(sessionCookieLang) == 0 {
		cookieLang, _ := c.Cookie("lang")
		headerLang := c.GetHeader("Accept-Language")

		session := sessions.Default(c)

		session.Set(constants.SessionCookieLang, cookieLang)
		session.Set(constants.SessionHeaderLang, headerLang)
		session.Save()
		sessionCookieLang = cookieLang
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"hideNav":                 true,
		constants.TemplateLangStr: sessionCookieLang,
	})
}

func Welcome(c *gin.Context) {
	c.HTML(http.StatusOK, "welcome.html", gin.H{
		constants.TemplateLangStr: c.GetString(constants.ContextSetLang),
	})
}
