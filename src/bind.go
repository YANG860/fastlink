package main

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserToken struct {
	jwt.RegisteredClaims
	ID int `json:"id"`
}

type LoginRequest struct {
	Account string `json:"account"`
	PW      string `json:"pw"`
}

type LoginResponse struct {
	Response
	Token string `json:"token"`
}
type Response struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type ShortUrlRequest struct {
	Token  string `json:"token"`
	Source string `json:"source"`
}

type ShortUrlResponse struct {
	Response
	Url string `json:"url"`
}

type RemoveLinkRequest struct {
	Token string `json:"token"`
	Url   string `json:"link"`
}

type RemoveLinkResponse struct {
	Response
}

type RemoveUserRequest struct {
	Token string `json:"token"`
}

type RemoveUserResponse struct {
	Response
}
