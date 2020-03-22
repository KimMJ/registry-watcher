package client

import (
	"encoding/json"
	"fmt"
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

func (c *client) Head(url string) (http.Header, error) {
	return c.client.Head(url)
}

func (c *client) getAPIEndpoint(endpoint string, repository string) (string, error) {
	req, err := http.NewRequest("HEAD", endpoint+"/v2/", nil)
	if err != nil {
		log.Error(err)
		return "", err
	}

	resp, err := c.DoReturnResponse(req)
	if err != nil {
		log.Error(err)
		return "", err
	}

	auth := resp.Header.Get("Www-Authenticate")
	splited := strings.Split(auth, "\"")
	realm := splited[1]
	service := splited[3]
	return fmt.Sprintf("%s?service=%s&scope=repository:%s:pull", realm, service, repository), nil
}

func (c *client) GetToken(endpoint string, username string, passwd string, repository string, insecureRegistry bool) (string, error) {

	var token models.Token

	if !strings.Contains(endpoint, "://") {
		if insecureRegistry {
			endpoint = "http://" + endpoint
		} else {
			endpoint = "https://" + endpoint
		}
	}

	url, err := c.getAPIEndpoint(endpoint, repository)
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.WithField("url", url).Debug("GetToken url")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err)
		return "", err
	}

	req.SetBasicAuth(username, passwd)

	data, err := c.client.Do(req)
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.WithFields(log.Fields{
		"json": string(data),
	}).Debug("response token")

	err = json.Unmarshal(data, &token)
	if err != nil {
		return "", err
	}

	curToken := token.GetToken()

	return curToken, nil
}

func (c *client) GetTag(endpoint, repository, token string, insecureRegistry bool) ([]byte, error) {
	if !strings.Contains(endpoint, "://") {
		if insecureRegistry {
			endpoint = "http://" + endpoint
		} else {
			endpoint = "https://" + endpoint
		}
	}

	url := endpoint + "/v2/" + repository + "/tags/list"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	return c.client.Do(req)
}

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
