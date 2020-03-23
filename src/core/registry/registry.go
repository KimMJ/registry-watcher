package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/kimmj/registry-watcher/src/common/utils"
	"github.com/kimmj/registry-watcher/src/core/notification/webhook"

	//"log"
	"strings"

	//commonHttp "github.com/kimmj/registry-watcher/src/common/http"
	"github.com/kimmj/registry-watcher/src/common/models"
	"github.com/kimmj/registry-watcher/src/core/registry/client"
	log "github.com/sirupsen/logrus"
)

// func GetToken(r.Endpoint string, r.Username string, r.Password string, repository string) {

// }
type ImageManifests map[string]models.ImageManifest

func getCreationDate(registryURL, token, repository, tag string) (time.Time, error) {
	url := fmt.Sprintf("%s/v2/%s/manifests/%s", registryURL, repository, tag)
	req, err := http.NewRequest("GET", url, nil)
	//req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v1+json")
	// req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Error(err)
		return time.Time{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	//req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	c := client.NewClient()
	resp, err := c.Do(req)
	if err != nil {
		log.Error(err)
		return time.Time{}, err
	}

	var dockerManifest models.DockerManifest
	err = json.Unmarshal(resp, &dockerManifest)
	if err != nil {
		log.Error(err)
		return time.Time{}, err
	}

	return dockerManifest.GetCreationDate(), nil
}

func getDigest(registryURL, token, repository, tag string) (string, error) {
	url := fmt.Sprintf("%s/v2/%s/manifests/%s", registryURL, repository, tag)
	req, err := http.NewRequest("GET", url, nil)
	// req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Error(err)
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	//req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	c := client.NewClient()
	var digest string
	resp, err := c.DoReturnResponse(req)
	if err != nil {
		log.Error(err)
		return "", err
	}

	//get header
	//defer resp.Body.Close()
	//data, err = ioutil.ReadAll(resp.Body)
	// bodyString := string(data)
	// fmt.Println(bodyString)
	digest = resp.Header.Get("Docker-Content-Digest")
	return digest, nil
	//var digest string
	//header, err := client.Head(url, &digest)
	//if err != nil {
	//	log.Error(err)
	//}
	//fmt.Println(header)
	//fmt.Println(digest)
}

func PollImage(registries models.Registries, webhookURL string) {
	var artifact models.Artifact
	for _, r := range registries.DockerRegistry {
		fmt.Printf("polling: %+v\n", r)
		endpoint := r.Endpoint

		if !strings.Contains(endpoint, "://") {
			if r.InsecureRegistry {
				endpoint = "http://" + endpoint
			} else {
				endpoint = "https://" + endpoint
			}
		}
		log.WithField("endpoint", endpoint).Debug("polling image")
		//type empty {}

		//sem := make(chan empty, N)
		c := client.NewClient()
		for _, image := range r.Images {
			repository := image

			log.WithFields(log.Fields{
				"json": r,
			}).Debug("poll Image")
			token, err := c.GetToken(endpoint, r.Username, r.Password, repository, r.InsecureRegistry)
			if err != nil {
				log.Error(err)
			}

			log.Debug("get token: ", token)

			var tagList models.TagList
			data, err := c.GetTag(endpoint, repository, token, r.InsecureRegistry)
			if err != nil {
				log.Error(err)
			}

			err = json.Unmarshal(data, &tagList)
			if err != nil {
				log.Error(err)
			}

			log.WithFields(log.Fields{
				"tags": tagList.Tags,
			}).Debug("got tags")

			imageManifests := ImageManifests{}

			for _, tag := range tagList.Tags {
				digest, err := getDigest(endpoint, token, repository, tag)
				creationDate, err := getCreationDate(endpoint, token, repository, tag)
				if err != nil {
					log.Error(err)
					continue
				}
				log.WithFields(log.Fields{
					"endpoint":   r.Endpoint,
					"repository": repository,
					"tag":        tag,
					"digest":     digest,
				}).Debug("got digest")

				imageManifest := models.ImageManifest{tag, digest, creationDate}
				id := hash(tag, digest)
				imageManifests[id] = imageManifest
			}

			compareJSON(r.Endpoint, image, &imageManifests, webhookURL, &artifact)
			writeJSON(&imageManifests, r.Endpoint, image)
		}
	}

	if len(artifact.Artifacts) > 0 {
		webhook.Send(webhookURL, artifact)
	}
}

func compareJSON(endpoint, image string, compare *ImageManifests, webhookURL string, artifact *models.Artifact) {
	imageManifests := ImageManifests{}
	readJSON(endpoint, image, &imageManifests)

	// var artifact models.Artifact

	created := time.Time{}
	var dockerArtifact models.DockerArtifact
	for k, v := range *compare {
		if _, ok := imageManifests[k]; !ok {
			compareManifest := (*compare)[k]
			log.WithFields(log.Fields{
				"key":   k,
				"value": v,
			}).Debug("find mismatch")
			log.WithFields(log.Fields{
				"oldCreate": created,
				"newCreate": compareManifest.CreationDate,
			}).Debug("CreationDate")

			if created.Sub(compareManifest.CreationDate).Microseconds() >= 0 {
				log.Debug("it is not recent image")
				continue
			}
			created = compareManifest.CreationDate
			log.WithFields(log.Fields{
				"newCreate": created,
			}).Debug("date updated")

			dockerArtifact = models.DockerArtifact{
				CustomKind: false,
				Reference:  endpoint + "/" + image + ":" + compareManifest.Tag,
				Name:       endpoint + "/" + image,
				Type:       "docker/image",
				Version:    compareManifest.Tag,
			}

		}
	}
	tmp := time.Time{}
	if created != tmp {
		(*artifact).AddItem(dockerArtifact)
	}
}

func readJSON(endpoint, image string, manifests *ImageManifests) {
	splited := strings.Split(image, "/")
	image = splited[len(splited)-1]
	dir := splited[:len(splited)-1]
	directory := fmt.Sprintf("db/%s/%s", endpoint, strings.Join(dir, "/"))
	filePath := fmt.Sprintf("%s/%s.json", directory, image)

	jsonFile, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Error(err)
	}

	//var imageManifests ImageManifests

	err = json.Unmarshal(jsonFile, manifests)
	if err != nil {
		log.Error(err)
	}

	log.WithFields(log.Fields{
		"json": *manifests,
	}).Debug("read JSON file")
}

func writeJSON(imageManifests *ImageManifests, endpoint string, image string) {
	prettyJSON := utils.PrettyPrintJSON(*imageManifests)

	splited := strings.Split(image, "/")
	image = splited[len(splited)-1]
	dir := splited[:len(splited)-1]

	directory := fmt.Sprintf("db/%s/%s", endpoint, strings.Join(dir, "/"))
	err := os.MkdirAll(directory, os.ModePerm)

	if err != nil {
		log.Error(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", directory, image), []byte(prettyJSON), 0644)
	if err != nil {
		log.Error(err)
	}
}

func hash(str ...string) string {
	ret := str[0]
	for _, s := range str[1:] {
		ret += "-" + s
	}
	return ret
}
