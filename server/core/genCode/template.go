package genCode

import "fmt"

var tpl_model = fmt.Sprintf(`package {{.LastModelPackage}}

// {{.StructName}} {{.Comment}}
type {{.StructName}} struct{
{{range .Columns}}		{{.MapperName}}	{{.Type}}	%s{{.Tag}} json:"{{.Name}}"{{if gt (len .ZhName) 0}} comment:"{{.ZhName}}"{{end}}%s {{if gt (len .Comment) 0}}// {{.Comment}}{{end}}
{{end}}	}


func ({{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
`, "`", "`")

var tpl_modelReq = fmt.Sprintf(`package {{if .Module}}{{.Module}}Req{{else}}req{{end}}

import "{{.GoMod}}/model/common/request"

// {{.StructName}} {{.Comment}}
type {{.StructName}}ListReq struct{
{{range .Columns}}{{if eq .SearchType "between"}} Start{{.MapperName}}	{{.Type}}	%sjson:"start_{{.Name}}"%s{{printf "\n"}} End{{.MapperName}}	{{.Type}}	%sjson:"end_{{.Name}}"%s{{printf "\n"}}{{end}}{{if eq .SearchType "like" "="}} {{.MapperName}}	{{.Type}} %sjson:"{{.Name}}"%s{{printf "\n"}}{{end}}{{end}} request.PageInfo
}



`, "`", "`", "`", "`", "`", "`")

var tpl_Repo = fmt.Sprintf(`package {{if .Module}}{{.Module}}Dao{{else}}dao{{end}}

import (
	{{if .Module}}"{{.GoMod}}/dao"{{end}}
	"{{.GoMod}}/{{.ModelPackage}}"
	"{{.GoMod}}/{{.RequestPackage}}"
	"context"
	"errors"
	"gorm.io/gorm"
)

type {{.StructName}}Repo struct {
	data *{{if .Module}}dao.{{end}}Data
}

func New{{.StructName}}Repo(data *{{if .Module}}dao.{{end}}Data) *{{.StructName}}Repo {
	repo := {{.StructName}}Repo{}
	repo.data = data
	return &repo
}

// Create{{.StructName}} 创建{{.TableZhName}}
func (repo *{{.StructName}}Repo) Create{{.StructName}}(ctx context.Context, data *{{.LastModelPackage}}.{{.StructName}}) error {
	err := repo.data.GetDB().WithContext(ctx).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *{{.StructName}}Repo) Find{{.StructName}}ById(ctx context.Context, id {{.PrimaryKeyType}}) (*{{.LastModelPackage}}.{{.StructName}}, error) {
	var resp {{.LastModelPackage}}.{{.StructName}}
	err := repo.data.GetDB().WithContext(ctx).First(&resp, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errcode.DBDataNotFound
	}
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (repo *{{.StructName}}Repo) Update{{.StructName}}(ctx context.Context, data *{{.LastModelPackage}}.{{.StructName}}) error {
	return repo.data.GetDB().WithContext(ctx).Save(data).Error
}

func (repo *{{.StructName}}Repo) Delete{{.StructName}}(ctx context.Context, id {{.PrimaryKeyType}}) error {
	return repo.data.GetDB().WithContext(ctx).Delete({{.LastModelPackage}}.{{.StructName}}{},id).Error
}

func (repo *{{.StructName}}Repo) {{.StructName}}List(ctx context.Context,in {{.LastRequestPackage}}.{{.StructName}}ListReq) (list []*{{.LastModelPackage}}.{{.StructName}},count int64, err error) {
	list = make([]*{{.LastModelPackage}}.{{.StructName}}, 0)
	tx:=repo.data.GetDB().WithContext(ctx)
	{{range .Columns}}{{if eq .SearchType "between"}}
	if in.Start{{.MapperName}} != "" {
		tx = tx.Where("{{.Name}} > ?",in.Start{{.MapperName}})
	}
	if in.End{{.MapperName}} != "" {
		tx = tx.Where("{{.Name}} < ?",in.Start{{.MapperName}})
	}
{{end}}{{if eq .SearchType "like"}}
	if {{.MapperName}} != "" {
		tx = tx.Where("{{.Name}} like ?","%%"+in.{{.MapperName}}+"%%")
	}
{{end}}{{if eq .SearchType "="}}
	if {{.MapperName}} != "" {
		tx = tx.Where("{{.Name}} = ?",in.{{.MapperName}})
	}
{{end}}{{end}}

	tx:= .Limit(in.Limit()).Offset(in.Offset())
	err = tx.Find(&list).Error
	if err != nil {
		return
	}
	err = tx.Count(&count).Error
	return
}




`)

var tpl_service = fmt.Sprintf(`package service
import (
 	"context"
	"{{.GoMod}}/model/{{.Module}}"
	"{{.GoMod}}/internal/dao/{{.Module}}"
)

type {{.StructName}}Service struct {
	{{.DownLatterStructName}}Repo *{{.StructName}}Repo
}

func New{{.StructName}}Service({{.DownLatterStructName}}Repo *{{.StructName}}Repo) *{{.StructName}}Service {
	service := {{.StructName}}Service{}
	service.{{.DownLatterStructName}}Repo = {{.StructName}}Repo
	return &service
}

// Create{{.StructName}} 创建{{.TableZhName}}
func (svc *{{.StructName}}Service) Create{{.StructName}}(ctx context.Context, data *model.{{.StructName}}) error {
	return svc.{{.DownLatterStructName}}Repo.Create{{.StructName}}(ctx, data)
}

func (svc *{{.StructName}}Service) Find{{.StructName}}(ctx context.Context, id {{.PrimaryKeyType}}) (*model.{{.StructName}}, error) {
	return svc.{{.DownLatterStructName}}Repo.FindOne(ctx, id)
}

func (svc *{{.StructName}}Service) Update{{.StructName}}(ctx context.Context, data *model.{{.StructName}}) error {
	return svc.{{.DownLatterStructName}}Repo.Update(ctx, data)
}

func (svc *{{.StructName}}Service) Delete{{.StructName}}(ctx context.Context, id {{.PrimaryKeyType}}) error {
	return svc.{{.DownLatterStructName}}Repo.Delete(ctx, id)
}

func (svc *{{.StructName}}Service) {{.StructName}}List(ctx context.Context,in request.{{.StructName}}Req) (*response.CountListResponse, error) {
	list,count,err:= svc.{{.DownLatterStructName}}Repo.List(ctx, in)
	if err != nil {
		return nil,err
	}
	return &response.CountListResponse{Count: count, List: list},nil
}
`)

var tpl_api = fmt.Sprintf(`package api
import (
	"{{.ServicePackage}}"
	"{{.ModelPackage}}"
	"{{.RequestPackage}}"
)

type {{.StructName}}Rest struct {
	{{.DownLatterStructName}}Service *{{.StructName}}Service
}

func New{{.StructName}}Rest({{.DownLatterStructName}}Service *{{.StructName}}Service) *{{.StructName}}Rest {
	rest := {{.StructName}}Rest{}
	rest.{{.DownLatterStructName}}Service = {{.StructName}}Service
	return &rest
}

func (rest {{.StructName}}Rest) Router(route *gin.RouterGroup) {
	group:=route.Group("{{.DownLatterStructName}}")
	group.POST("create{{.StructName}}",rest.create{{.StructName}})
	group.POST("update{{.StructName}}",rest.update{{.StructName}})
	group.POST("delete{{.StructName}}",rest.delete{{.StructName}})
}

//	create{{.StructName}} 创建{{.TableZhName}}
func (rest {{.StructName}}Rest) create{{.StructName}}(ctx *gin.Context) {
	var param {{.LastModelPackage}}.{{.StructName}}
	if err := ctx.ShouldBind(&param); err != nil {
		rest.Error(ctx, err)
		return
	}
	err = rest.{{.DownLatterStructName}}Service.Create{{.StructName}}(ctx, param)
	rest.Error(ctx, err)
}
//	update{{.StructName}} 修改{{.TableZhName}}
func (rest {{.StructName}}Rest) update{{.StructName}}(ctx *gin.Context) {
	var param {{.LastModelPackage}}.{{.StructName}}
	if err := ctx.ShouldBind(&param); err != nil {
		rest.Error(ctx, err)
		return
	}
	err = rest.{{.DownLatterStructName}}Service.Update{{.StructName}}(ctx, param)
	rest.Error(ctx, err)
}
//	delete{{.StructName}} 删除{{.TableZhName}}
func (rest {{.StructName}}Rest) delete{{.StructName}}(ctx *gin.Context) {
	var param request.{{.StructName}}
	if err := ctx.ShouldBind(&param); err != nil {
		rest.Error(ctx, err)
		return
	}
	err = rest.{{.DownLatterStructName}}Service.Delete{{.StructName}}(ctx, param)
	rest.Error(ctx, err)
}

`)
