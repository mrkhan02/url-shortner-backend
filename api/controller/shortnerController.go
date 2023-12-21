package controller

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/mrkhan02/url-shortner-api/database"
	"github.com/mrkhan02/url-shortner-api/helper"
	"github.com/mrkhan02/url-shortner-api/models"
	"github.com/redis/go-redis/v9"
)

func ShortenURL() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		err := godotenv.Load(".env")
		var body models.Request

		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// implement rate limiting
		redisClient := database.CreateClient(1)
		defer redisClient.Close()

		val, err := redisClient.Get(database.Ctx, c.ClientIP()).Result()
		if err == redis.Nil {
			_ = redisClient.Set(database.Ctx, c.ClientIP(), 10, 30*time.Minute).Err()
		} else {
			valInt, _ := strconv.Atoi(val)
			if valInt <= 0 {
				limit, _ := redisClient.TTL(database.Ctx, c.ClientIP()).Result()
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Rate limit exceeded",
					"rate_limit_reset": limit / time.Nanosecond / time.Minute})
				return
			}
		}

		// check if the imput is an actual url

		if !govalidator.IsURL(body.URL) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
			return
		}
		// check for domain error

		if !helper.RemoveDomainError(body.URL) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "BassG You can't Hack :)"})
			return

		}

		// enforce https, ssl

		body.URL = helper.EnforceHTTP(body.URL)

		var id string

		if body.CustomShort == "" {
			id = uuid.New().String()[:6]
		} else {
			id = body.CustomShort
		}

		r := database.CreateClient(0)
		defer r.Close()

		val, _ = r.Get(database.Ctx, id).Result()
		if val != "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "URL custom sort is already in use"})
			return
		}

		if body.Expiry == 0 {
			body.Expiry = 24
		}
		err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to server"})
			return
		}
		var Response models.Response
		Response.URl = body.URL
		Response.CustomShort = ""
		Response.Expiry = body.Expiry
		Response.XRateRemaining = 10
		Response.XRateLimitReset = 30

		redisClient.Decr(database.Ctx, c.ClientIP())
		val, _ = redisClient.Get(database.Ctx, c.ClientIP()).Result()
		Response.XRateRemaining, _ = strconv.Atoi(val)
		ttl, _ := redisClient.TTL(database.Ctx, c.ClientIP()).Result()
		Response.XRateLimitReset = ttl / time.Nanosecond / time.Minute
		Response.CustomShort = os.Getenv("DOMAIN") + "/" + id

		c.JSON(http.StatusOK, Response)
	}
}
