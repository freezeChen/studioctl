package genCode

import (
	"fmt"
	"github.com/freezeChen/studioctl/core/xresult"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"path"
)

func NewServer(port string) {
	engine := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	engine.Use(cors.New(config))

	engine.GET("gen/tables", tables)
	engine.GET("gen/columns", tableColumns)
	engine.POST("gen/preview", preview)
	engine.POST("gen/download", download)

	engine.Run(port)
}

func tables(ctx *gin.Context) {
	tables, err := query.Tables()
	xresult.OK(ctx, tables, err)
}

func tableColumns(ctx *gin.Context) {
	table := ctx.Query("table")
	column, err := query.TableColumn(table)
	if err != nil {
		xresult.Err(ctx, err)
		return
	}
	xresult.OK(ctx, parseTableColumns(table, column), err)
}

func preview(ctx *gin.Context) {
	var param PreviewReq
	if err := ctx.ShouldBind(&param); err != nil {
		xresult.Err(ctx, err)
		return
	}

	tableMapper := getTableMapper(param)
	code, err := genCode(tableMapper)
	xresult.OK(ctx, code, err)
}

func download(ctx *gin.Context) {
	var param PreviewReq
	if err := ctx.ShouldBind(&param); err != nil {
		xresult.Err(ctx, err)
		return
	}

	tableMapper := getTableMapper(param)
	code, err := genCode(tableMapper)
	if err != nil {
		xresult.Err(ctx, err)
		return
	}

	for _, v := range code.Codes {
		fmt.Println(path.Join(".", v.Path, v.FileName))
		os.MkdirAll(path.Join(".", v.Path), os.ModePerm)
		file, err := os.Create(path.Join(".", v.Path, v.FileName))
		if err != nil {

		}
		file.WriteString(v.Code)
		file.Close()
	}

	xresult.OK(ctx, nil, nil)

}
