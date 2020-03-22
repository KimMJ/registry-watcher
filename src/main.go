package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kimmj/registry-watcher/src/core/registry"
)

func main() {
	r := gin.Default()

	// cr := cron.New(cron.WithSeconds())
	// cr.AddFunc("*/5 * * * * *", PollImage)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/cron", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "cron started!",
		})
		// cr.Start()
	})

	r.GET("/poll", func(c *gin.Context) {
		registry.PollImage("https://wonderland-laptop.com", "admin", "Harbor12345", "test/busybox")
		c.JSON(200, gin.H{
			"message": "polling success",
		})
	})

	r.GET("/webhook", func(c *gin.Context) {
		// WebhookSender()
		c.JSON(200, gin.H{
			"message": "webhook is sended",
		})
	})

	r.GET("/readjson", func(c *gin.Context) {
		// ReadJsonFile()
		c.JSON(200, gin.H{
			"message": "read json",
		})
	})

	r.GET("/writejson", func(c *gin.Context) {
		// WriteJsonFile()
		c.JSON(200, gin.H{
			"message": "write json",
		})
	})

	r.GET("/comparejson", func(c *gin.Context) {
		// CompareJsonFile()
		c.JSON(200, gin.H{
			"message": "compare json",
		})
	})
	r.Run(":8888") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
