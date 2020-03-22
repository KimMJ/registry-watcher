package client

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"

	commonHttp "github.com/kimmj/registry-watcher/src/common/http"
	"github.com/kimmj/registry-watcher/src/common/models"
)

// Client defines methods that registry should implement
type Client interface {
}

type client struct {
	// baseURL string
	client *commonHttp.Client
}

func NewClient() *client {
	client := &client{
		client: commonHttp.NewClient(),
	}
	return client
}

func (c *client) Head(url string, digest *string) (http.Header, error) {
	return c.client.Head(url, digest)
}

func (c *client) GetToken(registryURL string, username string, passwd string, repository string) (string, error) {

	var token models.Token

	if !strings.Contains(registryURL, "://") {
		registryURL = "http://" + registryURL
	}

	url := registryURL + "/service/token?service=harbor-registry&scope=repository:" + repository + ":pull,push"
	//TODO: Use Header realm
	log.WithField("url", url).Debug("GetToken url")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(username, passwd)

	data, err := c.client.Do(req)
	if err != nil {
		return "", err
	}

	log.WithFields(log.Fields{
		"json": string(data),
	}).Debug("response token")
	// fmt.Println(data)
	//bodyString := string(data)
	//fmt.Println(bodyString)
	json.Unmarshal(data, &token)
	curToken := token.GetToken()

	return curToken, nil
}

func (c *client) GetTag(registryURL, repository, token string) ([]byte, error) {
	if !strings.Contains(registryURL, "://") { //TODO: if insecure-registry
		registryURL = "http://" + registryURL
	}
	url := registryURL + "/v2/" + repository + "/tags/list"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

// func (c *Client) Head(url string, digest *string) (http.Header, error) {
// 	req, err := http.NewRequest("HEAD", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp, err := c.client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp.Header, nil
// }

// func (c *Client) Post(url string, v ...interface{}) error {
// 	var reader io.Reader
// 	if len(v) > 0 {
// 		if r, ok := v[0].(io.Reader); ok {
// 			reader = r
// 		} else {
// 			data, err := json.Marshal(v[0])
// 			if err != nil {
// 				return err
// 			}

// 			reader = bytes.NewReader(data)
// 		}
// 	}

// 	req, err := http.NewRequest("POST", url, reader)
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	_, err = c.do(req)
// 	return err
// }

// func (c *client) do(req *http.Request) ([]byte, error) {
// 	resp, err := c.client.do(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer resp.Body.Close()
// 	data, err := ioutil.ReadAll(resp.Body)
// 	// bodyString := string(data)
// 	// fmt.Println(bodyString)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, err
// 	}

// 	return data, nil
// }
func (c *client) Do(req *http.Request) ([]byte, error) {
	return c.client.Do(req)
}

func (c *client) DoReturnResponse(req *http.Request) (*http.Response, error) {
	resp, err := c.client.DoReturnResponse(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}
