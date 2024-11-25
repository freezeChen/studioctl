package genCode

import (
	"fmt"
	"github.com/freezeChen/studioctl/core/util"
	"github.com/freezeChen/studioctl/core/xresult"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"strings"
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

	engine.GET("setting/loadGoInfo", loadTarget)

	engine.Run(port)
}

//region 表信息

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
	out, err := parseTableColumns(table, column)

	xresult.OK(ctx, out, err)
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
		if v.Type == 1 {
			os.MkdirAll(param.GoOutDir+"/"+path.Join(".", v.Path), os.ModePerm)
			file, err := os.Create(param.GoOutDir + "/" + path.Join(".", v.Path, v.FileName))
			if err != nil {
				fmt.Println("文件生成失败", err.Error())
			}
			file.WriteString(v.Code)
			file.Close()
		} else {
			os.MkdirAll(param.JsOutDir+"/"+path.Join(".", v.Path), os.ModePerm)
			file, err := os.Create(param.JsOutDir + "/" + path.Join(".", v.Path, v.FileName))
			if err != nil {
				fmt.Println("文件生成失败", err.Error())
			}
			file.WriteString(v.Code)
			file.Close()
		}

	}

	xresult.OK(ctx, nil, nil)

}

//endregion

//region 设置

func loadTarget(ctx *gin.Context) {
	targetDir := ctx.Query("target")
	workDir, _ := os.Getwd()
	fmt.Println(workDir)
	os.Chdir(targetDir)

	workDir, _ = os.Getwd()
	fmt.Println(workDir)

	modulePath, err := util.GetGoModuleName()
	if err != nil {
		xresult.Err(ctx, err)
		return
	}
	modulePath = strings.ReplaceAll(modulePath, "\\", "/")
	xresult.OK(ctx, modulePath, err)

}

//endregion
