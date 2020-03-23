package webhook

import (
	"bytes"
	"encoding/json"

	"github.com/kimmj/registry-watcher/src/common/http"
	"github.com/kimmj/registry-watcher/src/common/models"
	log "github.com/sirupsen/logrus"
)

func Send(targetURL string, artifact models.Artifact) {
	//busybox := DockerArtifact{false, "dockerrepo:8081/test/nginx:v2", "dockerrepo:8081/test/nginx", "docker/image", "v2"}
	// debian := DockerArtifact{false, "dockerrepo:8081/test/debian:v1", "dockerrepo:8081/test/debian", "docker/image", "v1"}

	c := http.NewClient()

	log.WithFields(log.Fields{
		"response": artifact,
	}).Debug("webhook data")

	// curl -i -X POST http://10.251.201.165:30200/webhooks/webhook/test --data @payload.json -H "Content-Type: application/json" --noproxy "*"

	//spinnakerUrl := "http://10.251.201.165:30200/webhooks/webhook/test"
	pbytes, _ := json.Marshal(artifact)
	// fmt.Println(pbytes)
	buff := bytes.NewBuffer(pbytes)
	err := c.Post(targetURL, buff)

	if err != nil {
		log.Error(err)
		return
	}

	//fmt.Println(buff)
	log.WithFields(log.Fields{
		"response": string(pbytes),
	}).Debug("webhook sended")

	//req, err := http.NewRequest("POST", spinnakerUrl, buff)
	//req.Header.Set("Content-Type", "application/json")
	//resp, err := c.Do(req)

	//defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//bodyString := string(body)
	//fmt.Println(bodyString)
}
