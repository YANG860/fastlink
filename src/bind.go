package main

import (
	"github.com/golang-jwt/jwt/v5"
)

type loginToken struct {
	Account string `json:"account"`
	PW      string `json:"pw"`
}

type userToken struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

type ShortUrl struct {
	Token  string `json:"token"`
	Source string `json:"source"`
}
