package main

import "github.com/gin-gonic/gin"

func removeUser(ctx *gin.Context) {

	var body RemoveUserRequest
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(400, RemoveUserResponse{Response: Response{Success: false, Error: "Invalid request"}})
		return

	}

	userToken, err := ParseJWT(body.Token)
	if err != nil {
		ctx.JSON(401, RemoveUserResponse{Response: Response{Success: false, Error: "Invalid token"}})
		return
	}

	_, err = engine.ID(userToken.ID).Cols("valid").Update(&User{Valid: false})
	if err != nil {
		ctx.JSON(500, RemoveUserResponse{Response: Response{Success: false, Error: "Database error"}})
		return
	}
	ctx.JSON(200, RemoveUserResponse{Response: Response{Success: true}})
}
