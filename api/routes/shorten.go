package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mrkhan02/url-shortner-api/controller"
)

func ShortRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/api/v1", controller.ShortenURL())
}
