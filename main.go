package main

import (
	"html/template"
	"log"
	"path"

	"github.com/aaa59891/go_demo_fullstack/configs"
	"github.com/aaa59891/go_demo_fullstack/controllers"
	"github.com/aaa59891/go_demo_fullstack/db"
	"github.com/aaa59891/go_demo_fullstack/inits"
	"github.com/aaa59891/go_demo_fullstack/routers"
	"github.com/aaa59891/go_demo_fullstack/templateFuncs"
	"github.com/aaa59891/go_demo_fullstack/utils"
	"github.com/boj/redistore"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gorillaSession "github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/i18n"
)

func init() {
	inits.CreateTable()
	inits.RegisterStruct()
}

func main() {
	defer db.DB.Close()
	dir := utils.GetProjectRoot()
	config := configs.GetConfig()
	r := gin.Default()

	i18n.MustLoadTranslationFile(path.Join(dir, "i18n", "en-us.all.json"))
	i18n.MustLoadTranslationFile(path.Join(dir, "i18n", "zh-tw.all.json"))
	r.Static("/static", path.Join(dir, "static"))
	r.SetFuncMap(template.FuncMap{
		"last":       templateFuncs.LastIndex,
		"dateTime24": templateFuncs.DateTime24,
		"date":       templateFuncs.Date,
		"actionName": templateFuncs.ActionName,
		"translate":  templateFuncs.I18nTranslate,
		"getPageLi":  templateFuncs.GetPageLi,
	})

	r.LoadHTMLGlob(path.Join(dir, "views/**/*"))

	store, err := redistore.NewRediStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		log.Fatal(err)
	}

	store.SetMaxLength(8092)

	r.Use(sessions.Sessions("mysession", &MyRedisStore{store}))
	r.GET("/socket.io/", gin.WrapH(inits.SocketServer))
	r.NoRoute(controllers.ErrorHandler404)
	routers.SetRoutes(r)

	r.Run(config.Server.Port)
}

/**
gin.session copy gorilla's session but not implement all the functionality,
so this struct does the same thing as the gin.session(wrap redistore.Redistore as gin.session's store)
*/
type MyRedisStore struct {
	*redistore.RediStore
}

func (c *MyRedisStore) Options(options sessions.Options) {
	c.RediStore.Options = &gorillaSession.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
