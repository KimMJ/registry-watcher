package models

import (
	"time"
)

type Token struct {
	Token       string    `json:"token"`
	AccessToken string    `json:"access_token"` // used in azure container
	ExpiresIn   int       `json:"expires_in"`
	IssuedAt    time.Time `json:"issued_at"`
}

func (t *Token) GetToken() string {
	token := t.Token
	if len(token) == 0 {
		token = t.AccessToken
	}
	return token
}
