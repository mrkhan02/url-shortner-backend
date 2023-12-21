package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mrkhan02/url-shortner-api/controller"
)

func ResRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/:url", controller.ResolveURL())
}
