package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserToken struct {
	jwt.RegisteredClaims
	ID      int    `json:"id"`
	Account string `json:"account"`
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

type GetUserRequest struct {
	Token string `json:"token"`
}

type GetUserResponse struct {
	Response
	Account         string    `json:"account"`
	LinkCount       int       `json:"link_count"`
	RegisteredAt    time.Time `json:"registered_at"`
	LatestCreatedAt time.Time `json:"latest_created_at"`
}

type GetLinkRequest struct {
	Token string `json:"token"`
}

type GetLinkResponse struct {
	Response
	SourceUrl  string    `json:"source"`
	ShortUrl   string    `json:"url"`
	CreatedAt  time.Time `json:"created_at"`
	ExpireAt   time.Time `json:"expire_at"`
	ClickCount int       `json:"click_count"`
}
