package basic

import (
	"net/http"
)

type authorizer struct {
	username string
	password string
}

func NewAuthorizer(username, password string) *authorizer {
	return &authorizer{
		username: username,
		password: password,
	}
}

func (a *authorizer) Modify(req *http.Request) error {
	if len(a.username) > 0 {
		req.SetBasicAuth(a.username, a.password)
	}
	return nil
}
