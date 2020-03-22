package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kimmj/registry-watcher/src/common/models"
	"github.com/kimmj/registry-watcher/src/core/registry"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	cr *cron.Cron
)

func init() {
	// set log level
	log.SetLevel(log.DebugLevel)

	// read config file
	config := models.Config{}
	err := config.ReadConfig("src/config.yml")
	if err != nil {
		log.Println("Read config file got err")
	}

	log.WithFields(log.Fields{
		"config": fmt.Sprintf("%+v", config),
	}).Debug("load config")

	cr = cron.New(cron.WithSeconds())

	for _, webhook := range config.Webhook {
		for _, dockerRegistry := range webhook.Registries.DockerRegistry {
			_, err := cr.AddFunc("*/5 * * * * *", func() {
				registry.PollImage(&dockerRegistry, webhook.EndPoint)
			})
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/cron", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "cron started!",
		})
		cr.Start()
	})

	// for Test
	r.GET("/poll", func(c *gin.Context) {
		dockerRegistry := models.DockerRegistry{"wonderland-laptop.com", "admin", "Harbor12345", false, []string{"test/busybox"}}
		registry.PollImage(&dockerRegistry, "http://192.168.8.22:30200/webhooks/webhook/test")
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
	err := r.Run(":8888") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		log.Error(err)
	}
}
