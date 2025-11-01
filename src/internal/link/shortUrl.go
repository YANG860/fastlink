package link

import (
	"fastlink/auth"
	"fastlink/db"
	"fastlink/models"
	"math/rand"
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
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// GetShortUrl 生成短链，需验证 token 和原始链接有效性
func GetShortUrl(ctx *gin.Context) {

	var body models.ShortUrlRequest
	// 解析请求体
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(400, models.InvalidRequestError)
		return
	}

	// 校验 token
	token, err := auth.ParseJWT(body.Token)

	if err != nil {
		ctx.JSON(401, models.InvalidTokenError)
		return
	}

	// 校验原始链接格式
	if !checkUrl(body.Source) {
		if !checkUrl("https://" + body.Source) {
			ctx.JSON(400, models.InvalidRequestError)
			return
		}
		body.Source = "https://" + body.Source
	}
	// 生成唯一短链
	s := genRandomString(6)
	for has, err := db.SQLEngine.Exist(&db.Link{ShortUrl: s}); has; {
		if err != nil {
			ctx.JSON(500, models.DatabaseError)
			return
		}
		s = genRandomString(6)
	}

	var user db.User
	db.SQLEngine.ID(token.ID).Get(&user)
	if !user.Valid {
		ctx.JSON(401, models.InvalidTokenError)
		return
	}

	// 事务：插入短链并更新用户信息
	_, err = db.SQLEngine.Transaction(func(tx *xorm.Session) (interface{}, error) {

		_, err := tx.InsertOne(&db.Link{
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
		ctx.JSON(500, models.DatabaseError)
		return
	}

	ctx.JSON(200, models.ShortUrlResponse{
		Response: models.Success,
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
