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
)

type Harbor_Response struct {
    Token            string
    Access_token     string
    Expires_in       int
    Issued_at        string
}

type TagList struct {
    Name    string
    Tags    []string
}

func tt() {
    fmt.Println(time.Now())
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
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
        tr := &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        client := &http.Client{Transport: tr}
        var harbor string = "https://wonderland-laptop.com"
        var username string = "admin"
        var passwd string = "Harbor12345"
        var repository string = "test/busybox"
        var url string = harbor + "/service/token?service=harbor-registry&scope=repository:" + repository + ":pull,push"
        fmt.Println(url)

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

        c.JSON(200, gin.H{
            "message": "polling success",
        })
    })
  r.Run(":12345") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
