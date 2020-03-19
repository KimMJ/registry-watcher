package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "crypto/tls"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "github.com/robfig/cron/v3"
    "time"
    "bytes"
    "net"
)

type Harbor_Response struct {
    Token            string     `json:"token"`
    Access_token     string     `json:"access_token"`
    Expires_in       int        `json:"expires_in"`
    Issued_at        string     `json:"issued_at"`
}

type TagList struct {
    Name            string      `json:"name"`
    Tags            []string    `json:"tags"`
}

type DockerArtifact struct {
    CustomKind      bool        `json:"customKind"`
    Reference       string      `json:"reference"`
    Name            string      `json:"name"`
    Type            string      `json:"type"`
    Version         string      `json:"version"`
}

type Artifact struct {
    Artifacts    []DockerArtifact `json:"artifacts"`
}

func tt() {
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
    req.Header.Set("Authorization", "Bearer " + curToken)
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
    req.Header.Set("Authorization", "Bearer " + curToken)
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

func (artifacts *Artifact) AddItem(item DockerArtifact) []DockerArtifact {
    artifacts.Artifacts = append(artifacts.Artifacts, item)
    return artifacts.Artifacts
}

func main() {
    r := gin.Default()

    cr := cron.New(cron.WithSeconds())
    cr.AddFunc("*/5 * * * * *", tt)

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
        tt()
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
  r.Run(":12345") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
