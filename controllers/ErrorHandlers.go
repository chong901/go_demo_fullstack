package controllers

import (
	"github.com/aaa59891/mosi_demo_go/constants"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandler404(c *gin.Context) {
	GoToErrorPage(http.StatusNotFound, c, constants.Err404NotFounr)
}
