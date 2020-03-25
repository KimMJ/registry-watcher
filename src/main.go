package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kimmj/registry-watcher/src/common/http"
	"github.com/kimmj/registry-watcher/src/common/models"
	"github.com/kimmj/registry-watcher/src/core/registry"
	cron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	cr           *cron.Cron
	config       models.Config
	test         models.Registries
	testEndpoint string
)

func init() {
	// set log level
	log.SetLevel(log.DebugLevel)

	// read config file
	config = models.Config{}
	err := config.ReadConfig("src/config.yml")
	if err != nil {
		log.Println("Read config file got err")
	}

	log.WithFields(log.Fields{
		"config": fmt.Sprintf("%+v", config),
	}).Debug("load config")

	cr = cron.New(cron.WithSeconds())
	for j, webhook := range config.Webhook {
		// for i, dockerRegistry := range webhook.Registries.DockerRegistry {
		log.WithFields(log.Fields{
			"function":       fmt.Sprintf("registry.PollImage(&dockerRegistry, %s)", webhook.EndPoint),
			"dockerRegistry": fmt.Sprintf("%+v", config.Webhook[j].Registries),
		}).Debug("cron added")

		test = config.Webhook[j].Registries
		testEndpoint = webhook.EndPoint
		_, err := cr.AddFunc("*/30 * * * * *", func() {
			registry.PollImage(test, webhook.EndPoint)
		})

		if err != nil {
			log.Error(err)
		}

		// }
	}
	http.InitSema()
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
		registry.PollImage(test, testEndpoint)
		c.JSON(200, gin.H{
			"message": "polling success",
		})
	})

	err := r.Run(":8888") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		log.Error(err)
	}
}
