package registry

import (
	"encoding/json"
	"fmt"

	// "github.com/kimmj/registry-watcher/src/common/http"
	"github.com/kimmj/registry-watcher/src/common/models"
	"github.com/kimmj/registry-watcher/src/core/registry/client"
)

// func GetToken(registryURL string, username string, passwd string, repository string) {

// }

func getTag(registryURL, repository, token string) {

}

func PollImage(registryURL string, username string, passwd string, repository string) {
	// fmt.Println(time.Now())
	c := client.NewClient()
	token, err := c.GetToken(registryURL, username, passwd, repository)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(token)

	var tagList models.TagList
	data, err := c.GetTag(registryURL, repository, token)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
	err = json.Unmarshal(data, &tagList)
	fmt.Println(tagList.Tags)

	// url = harbor + "/v2/" + repository + "/tags/list"
	// var tagList models.TagList
	// err = client.Get(url, &tagList)

	// req, err = http.NewRequest("GET", url, nil)
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer "+curToken)
	// resp, err = client.Do(req)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// defer resp.Body.Close()
	// data, err = ioutil.ReadAll(resp.Body)
	// // bodyString = string(data)
	// // fmt.Println(bodyString)
	// var tagList TagList
	// err = json.Unmarshal(data, &tagList)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(tagList.Tags)

	// url = harbor + "/v2/" + repository + "/manifests/v1"

	// req, err = http.NewRequest("GET", url, nil)
	// // req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer "+curToken)
	// req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	// resp, err = client.Do(req)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// //get header
	// defer resp.Body.Close()
	// data, err = ioutil.ReadAll(resp.Body)
	// // bodyString := string(data)
	// // fmt.Println(bodyString)
	// fmt.Println(resp.Header.Get("Docker-Content-Digest"))

	// var digest string
	// header, err := client.Head(url, &digest)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(header)
	// fmt.Println(digest)
}
