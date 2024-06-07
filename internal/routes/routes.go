package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/niazlv/sport-plus-LCT/internal/api/auth"
)

func Setup(app *gin.Engine) {
	api := app.Group("v1")

	auth.Setup(api)
}
