package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	client   *http.Client
	username string
	password string
}

const (
	semaLimit = 500
)

var Sema chan struct{}

func InitSema() {
	Sema = make(chan struct{}, semaLimit)
}

func NewClient() *Client {
	client := &Client{}

	var transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client.client = &http.Client{Transport: transport}

	return client
}

func (c *Client) Get(url string, v ...interface{}) error {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Error(err)
		return err
	}

	data, err := c.do(req)
	if err != nil {
		log.Error(err)
		return err
	}

	if len(v) == 0 {
		log.Error(err)
		return nil
	}

	return json.Unmarshal(data, v[0])
}

func (c *Client) Head(url string) (http.Header, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return resp.Header, nil
}

func (c *Client) Post(url string, v ...interface{}) error {
	var reader io.Reader
	if len(v) > 0 {
		if r, ok := v[0].(io.Reader); ok {
			reader = r
		} else {
			data, err := json.Marshal(v[0])
			if err != nil {
				log.Error(err)

				return err
			}

			reader = bytes.NewReader(data)
		}
	}

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.do(req)
	log.WithFields(log.Fields{
		"resp": string(resp),
	}).Debug("post sended")
	return err
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	resp, err := c.DoReturnResponse(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if resp != nil {
		defer errCheckClose(resp)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	// bodyString := string(data)
	// fmt.Println(bodyString)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return data, nil
}

func (c *Client) Do(req *http.Request) ([]byte, error) {
	Sema <- struct{}{}
	ret, err := c.do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	<-Sema
	return ret, nil
}

func (c *Client) DoReturnResponse(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return resp, err
}

func errCheckClose(resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		log.Error(err)
	}
}
