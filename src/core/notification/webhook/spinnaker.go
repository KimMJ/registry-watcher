package webhook

import (
	"bytes"
	"encoding/json"

	"github.com/kimmj/registry-watcher/src/common/http"
	"github.com/kimmj/registry-watcher/src/common/models"
	log "github.com/sirupsen/logrus"
)

func Send(targetURL string, artifact models.Artifact) {
	c := http.NewClient()

	log.WithFields(log.Fields{
		"response": artifact,
	}).Debug("webhook data")

	pbytes, _ := json.Marshal(artifact)
	buff := bytes.NewBuffer(pbytes)
	err := c.Post(targetURL, buff)

	if err != nil {
		log.Error(err)
		return
	}

	log.WithFields(log.Fields{
		"response": string(pbytes),
	}).Debug("webhook sended")
}
