package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrkhan02/url-shortner-api/database"
	"github.com/redis/go-redis/v9"
)

func ResolveURL() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Param("url")

		r := database.CreateClient(0)
		defer r.Close()

		value, err := r.Get(database.Ctx, url).Result()

		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "short not found in the database"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot connect to DB"})
			return

		}

		rInr := database.CreateClient(1)
		defer rInr.Close()

		_ = rInr.Incr(database.Ctx, "counter")

		c.Redirect(http.StatusMovedPermanently, value)
		return

	}
}
