package models

import (
	"fastlink/db"
	"time"
)

// 通用响应结构体
type Response struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// 登录请求结构体
type LoginRequest struct {
	Account string `json:"account"`
	PW      string `json:"pw"`
}

// 登录响应结构体
type LoginResponse struct {
	Response
	Token string `json:"token"`
}

// 注册请求结构体
type RegisterRequest struct {
	Account string `json:"account"`
	PW      string `json:"pw"`
}

// 注册响应结构体
type RegisterResponse struct {
	Response
	Token string `json:"token"`
}

// 短链生成请求结构体
type ShortUrlRequest struct {
	Token  string `json:"token"`
	Source string `json:"source"`
}

// 短链生成响应结构体
type ShortUrlResponse struct {
	Response
	ShortUrl string `json:"short_url"`
}

// 删除短链请求结构体
type RemoveLinkRequest struct {
	Token string `json:"token"`
}

// 删除短链响应结构体
type RemoveLinkResponse struct {
	Response
}

// 删除用户请求结构体
type RemoveUserRequest struct {
	Token string `json:"token"`
}

// 删除用户响应结构体
type RemoveUserResponse struct {
	Response
}

// 获取用户信息请求结构体
type GetUserRequest struct {
	Token string `json:"token"`
}

// 获取用户信息响应结构体
type GetUserResponse struct {
	Response
	Account         string    `json:"account"`
	LinkCount       int       `json:"link_count"`
	RegisteredAt    time.Time `json:"registered_at"`
	LatestCreatedAt time.Time `json:"latest_created_at"`
}

// 获取短链信息请求结构体
type GetLinkRequest struct {
	Token string `json:"token"`
}

// 获取短链信息响应结构体
type GetLinkResponse struct {
	Response
	SourceUrl  string    `json:"source"`
	ShortUrl   string    `json:"url"`
	CreatedAt  time.Time `json:"created_at"`
	ExpireAt   time.Time `json:"expire_at"`
	ClickCount int       `json:"click_count"`
}

type Admin struct {
	Account string `json:"account"`
	PW      string `json:"pw"`
}

// 获取所有短链请求结构体
type GetAllLinksAdminRequest struct {
	Admin
}

// 获取所有短链响应结构体
type GetAllLinksAdminResponse struct {
	Response
	Links []db.Link `json:"links"`
}

// 删除短链请求结构体
type RemoveLinkAdminRequest struct {
	Admin
	ShortUrl string `json:"short_url"`
}

// 删除短链响应结构体
type RemoveLinkAdminResponse struct {
	Response
}

// 获取某用户所有短链请求结构体
type GetAllLinksByUserAdminRequest struct {
	Admin
	UserID int `json:"user_id"`
}

// 获取某用户所有短链响应结构体
type GetAllLinksByUserAdminResponse struct {
	Response
	Links []db.Link `json:"links"`
}

// 获取所有用户请求结构体
type GetAllUsersAdminRequest struct {
	Admin
}

// 获取所有用户响应结构体
type GetAllUsersAdminResponse struct {
	Response
	Users []db.User `json:"users"`
}

// 删除用户请求结构体
type RemoveUserAdminRequest struct {
	Admin
	UserID int `json:"user_id"`
}

// 删除用户响应结构体
type RemoveUserAdminResponse struct {
	Response
}
