package main

import (
	"math/rand/v2"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"xorm.io/xorm"
)

var letter string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890+="
var length int = 5

func getShortUrl(ctx *gin.Context) {

	s := strings.Builder{}
	for i := 0; i < length; i++ {
		s.WriteByte(letter[rand.IntN(len(letter))])
	}

	var body ShortUrl
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if body.Token == "" {
		cookieToken, err := ctx.Cookie("token")
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Failed to retrieve token cookie"})
			return
		}

		if cookieToken == "" {
			ctx.JSON(401, gin.H{"error": "No login"})
			return
		}

		body.Token = cookieToken
	}

	jwtSecret := []byte("key")
	rawToken, err := jwt.ParseWithClaims(body.Token, &userToken{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !rawToken.Valid {
		ctx.JSON(401, gin.H{"error": "Invalid token"})
		return
	}

	var token *userToken
	token, ok := rawToken.Claims.(*userToken)
	if !ok {
		ctx.JSON(401, gin.H{"error": "Invalid token claims"})
		return
	}

	if time.Now().After(token.ExpiresAt.Time) {
		ctx.JSON(401, gin.H{"error": "Invalid token"})
		return
	}

	var user User
	engine.ID(token.ID).Get(&user)
	if !user.Valid {
		ctx.JSON(401, gin.H{"error": "Invalid  user"})
		return
	}

	if !checkUrl(body.Source) {
		if !checkUrl("https://" + body.Source) {
			ctx.JSON(400, gin.H{"error": "Invalid  source url"})
		}
		body.Source = "https://" + body.Source
	}

	engine.Transaction(func(tx *xorm.Session) (interface{}, error) {

		user.LatestCreatedAt = time.Now()
		user.LinkCount++

		_, err := tx.InsertOne(&Link{
			SourceUrl: body.Source,
			ShortUrl:  s.String(),
			Userid:    user.ID,

			ExpireAt: time.Now().Add(time.Hour * 24 * 7),
		})

		if err != nil {
			return nil, err
		}


		_, err = tx.ID(user.ID).Cols("link_count", "latest_created_at").Update(&user)
		if err != nil {
			return nil, err
		}
		return nil, nil

	})

	ctx.JSON(200, gin.H{
		"success": true,
		"url":     "localhost:8080/" + s.String(),
	})

}

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
