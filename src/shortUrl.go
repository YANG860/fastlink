package main

import (
	"math/rand/v2"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

// genRandomString 生成指定长度的随机字符串
func genRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890+="
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.IntN(len(letters))]
	}
	return string(b)
}

// getShortUrl 生成短链，需验证 token 和原始链接有效性
func getShortUrl(ctx *gin.Context) {

	var body ShortUrlRequest
	// 解析请求体
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(400, Response{Success: false, Error: "Invalid request"})
		return
	}

	// 获取 token（优先请求体，否则 cookie）
	if body.Token == "" {
		cookieToken, err := ctx.Cookie("token")
		if err != nil {
			ctx.JSON(400, Response{Success: false, Error: "Failed to retrieve token cookie"})
			return
		}

		if cookieToken == "" {
			ctx.JSON(401, Response{Success: false, Error: "No login"})
			return
		}

		body.Token = cookieToken
	}

	// 校验 token
	token, err := ParseJWT(body.Token)

	if err != nil {
		ctx.JSON(401, Response{Success: false, Error: "Invalid token"})
		return
	}

	// 校验原始链接格式
	if !checkUrl(body.Source) {
		if !checkUrl("https://" + body.Source) {
			ctx.JSON(400, Response{Success: false, Error: "Invalid source url"})
			return
		}
		body.Source = "https://" + body.Source
	}
	// 生成唯一短链
	s := genRandomString(6)
	for has, err := engine.Exist(&Link{ShortUrl: s}); has; {
		if err != nil {
			ctx.JSON(500, Response{Success: false, Error: "Database error"})
			return
		}
		s = genRandomString(6)
	}

	var user User
	engine.ID(token.ID).Get(&user)
	if !user.Valid {
		ctx.JSON(401, Response{Success: false, Error: "Invalid user"})
		return
	}

	// 事务：插入短链并更新用户信息
	_, err = engine.Transaction(func(tx *xorm.Session) (interface{}, error) {

		_, err := tx.InsertOne(&Link{
			SourceUrl: body.Source,
			ShortUrl:  s,
			UserID:    user.ID,

			ExpireAt: time.Now().Add(time.Hour * 24 * 7),
		})

		if err != nil {
			return nil, err
		}

		user.LatestCreatedAt = time.Now()
		user.LinkCount++

		_, err = tx.ID(user.ID).Cols("link_count", "latest_created_at").Update(&user)
		if err != nil {
			return nil, err
		}
		return nil, nil

	})

	if err != nil {
		ctx.JSON(500, Response{Success: false, Error: "Database error"})
		return
	}

	ctx.JSON(200, ShortUrlResponse{
		Response: Response{
			Success: true,
		},
		ShortUrl: s,
	})

}

// checkUrl 校验 URL 格式是否合法
func checkUrl(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	return true
}
