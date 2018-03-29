package controllers

import (
	"net/http"

	"github.com/aaa59891/go_demo_fullstack/constants"
	"github.com/gin-gonic/gin"
)

func ErrorHandler404(c *gin.Context) {
	GoToErrorPage(http.StatusNotFound, c, constants.Err404NotFounr)
}
