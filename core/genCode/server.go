package genCode

import (
	"github.com/freezeChen/studioctl/core/xresult"
	"github.com/gin-gonic/gin"
)

func NewServer(port string) {
	engine := gin.Default()
	engine.GET("gen/tables", tables)
	engine.GET("gen/columns", tableColumns)
	engine.GET("gen/preview", preview)

	engine.Run(port)
}

func tables(ctx *gin.Context) {
	tables, err := query.Tables()
	xresult.OK(ctx, tables, err)
}

func tableColumns(ctx *gin.Context) {
	table := ctx.Query("table")
	column, err := query.TableColumn(table)
	xresult.OK(ctx, column, err)
}

func preview(ctx *gin.Context) {
	var param PreviewReq
	if err := ctx.ShouldBind(&param); err != nil {
		xresult.Err(ctx, err)
		return
	}

	model, err := genModel(param)

	xresult.OK(ctx, model, err)

}
