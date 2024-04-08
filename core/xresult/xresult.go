package xresult

import "github.com/gin-gonic/gin"

func OK(ctx *gin.Context, data any, err error) {
	if err != nil {
		ctx.String(500, err.Error())
	} else {
		ctx.JSON(200, data)
	}
}

func Err(ctx *gin.Context, err error) {
	ctx.String(500, err.Error())
}
