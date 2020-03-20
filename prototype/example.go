package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	// "gopkg.in/yaml.v2"
)

type Harbor_Response struct {
	Token        string `json:"token"`
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
	Issued_at    string `json:"issued_at"`
}

type TagList struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type DockerArtifact struct {
	CustomKind bool   `json:"customKind"`
	Reference  string `json:"reference"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Version    string `json:"version"`
}

type Artifact struct {
	Artifacts []DockerArtifact `json:"artifacts"`
}

type Manifest struct {
	Tag          string `json:"tag"`
	Digest       string `json:"digest"`
	CreationDate string `json:"creationDate"`
}

// type ImageManifests struct {
//     Images           []Manifest          `json:"images"`
// }

type ImageManifests map[string]Manifest

func PollImage() {
	fmt.Println(time.Now())
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{Transport: tr}
	//get token
	var harbor string = "https://wonderland-laptop.com"
	var username string = "admin"
	var passwd string = "Harbor12345"
	var repository string = "test/busybox"
	var url string = harbor + "/service/token?service=harbor-registry&scope=repository:" + repository + ":pull,push"
	// fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	// bodyString := string(data)
	// fmt.Println(bodyString)

	if err != nil {
		fmt.Println(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to start GC: %d", resp.StatusCode)
	}

	// fmt.Println("success")
	var harResponse Harbor_Response
	err = json.Unmarshal(data, &harResponse)
	if err != nil {
		fmt.Println(err)
	}
	curToken := harResponse.Token

	//get tags
	// curl -i -k -H "Content-Type: application/json" -H "Authorization:  Bearer token" -X GET https://wonderland-laptop/v2/test/busybox/tags/list
	url = harbor + "/v2/" + repository + "/tags/list"
	req, err = http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+curToken)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	// bodyString = string(data)
	// fmt.Println(bodyString)
	var tagList TagList
	err = json.Unmarshal(data, &tagList)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tagList.Tags)

	url = harbor + "/v2/" + repository + "/manifests/v1"
	req, err = http.NewRequest("GET", url, nil)
	// req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+curToken)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	//get header
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	// bodyString := string(data)
	// fmt.Println(bodyString)
	fmt.Println(resp.Header.Get("Docker-Content-Digest"))
}

func webhookSender() {
	busybox := DockerArtifact{false, "dockerrepo:8081/test/nginx:v2", "dockerrepo:8081/test/nginx", "docker/image", "v2"}
	// debian := DockerArtifact{false, "dockerrepo:8081/test/debian:v1", "dockerrepo:8081/test/debian", "docker/image", "v1"}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{Transport: tr}
	var data Artifact
	data.AddItem(busybox)
	// data.AddItem(debian)

	// fmt.Printf("%+v\n", data)
	// curl -i -X POST http://10.251.201.165:30200/webhooks/webhook/test --data @payload.json -H "Content-Type: application/json" --noproxy "*"

	spinnakerUrl := "http://10.251.201.165:30200/webhooks/webhook/test"
	pbytes, _ := json.Marshal(data)
	// fmt.Println(pbytes)
	buff := bytes.NewBuffer(pbytes)
	fmt.Println(buff)
	req, err := http.NewRequest("POST", spinnakerUrl, buff)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	fmt.Println(bodyString)

}

func ReadJsonFile() {
	registry := "wonderland-laptop.com"
	image := "test/busybox"
	jsonFile, err := ioutil.ReadFile("./db/" + registry + "/" + image + ".json")
	if err != nil {
		fmt.Println(err)
	}
	var imageManifests ImageManifests
	err = json.Unmarshal(jsonFile, &imageManifests)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", imageManifests)
}

func CompareJsonFile() {
	registry := "wonderland-laptop.com"
	image := "test/busybox"
	jsonFile, err := ioutil.ReadFile("./db/" + registry + "/" + image + ".json")
	if err != nil {
		fmt.Println(err)
	}
	var imageManifests ImageManifests
	err = json.Unmarshal(jsonFile, &imageManifests)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(PollImage())
}

func WriteJsonFile() {
	registry := "wonderland-laptop.com"
	image := "test/busybox"
	// jsonFile, err := ioutil.ReadFile("./db/" + registry + "/" + image + ".json")
	// if err != nil {
	//     fmt.Println(err)
	// }
	imageManifests := ImageManifests{}
	// var imageManifests ImageManifests
	// err = json.Unmarshal(jsonFile, &imageManifests)
	// if err != nil {
	//     fmt.Println(err)
	// }

	manifest := Manifest{"v2", "12345", "123"}
	manifestv1 := Manifest{"v1", "54321", "123"}
	fmt.Printf("%+v\n", manifest)

	//TODO: with hash
	id := manifest.Tag + "-" + manifest.Digest
	id2 := manifestv1.Tag + "-" + manifestv1.Digest
	// if imageManifests[id] == nil {
	// imageManifests[id] = make(map[string]Manifest)
	// }
	imageManifests[id] = manifest
	imageManifests[id2] = manifestv1
	fmt.Printf("%+v\n", imageManifests)
	d, err := json.MarshalIndent(&imageManifests, "", "\t")
	// d, err := json.MarshalIndent(map[string]interface{}{id: manifest}, "", "\t")

	// d, err := json.Marshal(&imageManifests)
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile("./db/"+registry+"/"+image+".json", d, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (artifacts *Artifact) AddItem(item DockerArtifact) []DockerArtifact {
	artifacts.Artifacts = append(artifacts.Artifacts, item)
	return artifacts.Artifacts
}

func main() {
	r := gin.Default()

	cr := cron.New(cron.WithSeconds())
	cr.AddFunc("*/5 * * * * *", PollImage)

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
		fmt.Println("hi!")
	})

	r.GET("/poll", func(c *gin.Context) {
		PollImage()
		c.JSON(200, gin.H{
			"message": "polling success",
		})
	})

	r.GET("/webhook", func(c *gin.Context) {
		webhookSender()
		c.JSON(200, gin.H{
			"message": "webhook is sended",
		})
	})

	r.GET("/readjson", func(c *gin.Context) {
		ReadJsonFile()
		c.JSON(200, gin.H{
			"message": "read json",
		})
	})

	r.GET("/writejson", func(c *gin.Context) {
		WriteJsonFile()
		c.JSON(200, gin.H{
			"message": "write json",
		})
	})

	r.GET("/comparejson", func(c *gin.Context) {
		CompareJsonFile()
		c.JSON(200, gin.H{
			"message": "write json",
		})
	})
	r.Run(":12345") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
