package genCode

import "fmt"

var tpl_model = fmt.Sprintf(`package model

// {{.StructName}} {{.Comment}}
type {{.StructName}} struct{
{{range .Columns}}		{{.MapperName}}	{{.Type}}	%s{{.Tag}} json:"{{.Name}}"{{if gt (len .ZhName) 0}} comment:"{{.ZhName}}"{{end}}%s {{if gt (len .Comment) 0}}// {{.Comment}}{{end}}
{{end}}	}


func ({{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
`, "`", "`")

var tpl_modelReq = fmt.Sprintf(`package request
// {{.StructName}} {{.Comment}}
type {{.StructName}}ListReq struct{
{{range .Columns}}{{if eq .SearchType "between"}} Start{{.MapperName}}	{{.Type}}	%sjson:"start_{{.Name}}"%s{{printf "\n"}} End{{.MapperName}}	{{.Type}}	%sjson:"end_{{.Name}}"%s{{printf "\n"}}{{end}}{{if eq .SearchType "like" "="}} {{.MapperName}}	{{.Type}} %sjson:"{{.Name}}"%s{{printf "\n"}}{{end}}{{end}} NumberPage
}



`, "`", "`", "`", "`", "`", "`")

var tpl_Repo = fmt.Sprintf(`package dao

import (
 	"context"
	"{{.GoMod}}/model"
)

type {{.StructName}}Repo struct {
	data *Data
}

func New{{.StructName}}Repo(data *Data) *{{.StructName}}Repo {
	repo := {{.StructName}}Repo{}
	repo.data = data
	return &repo
}

// Create{{.StructName}} 创建{{.TableZhName}}
func (repo *{{.StructName}}Repo) Create{{.StructName}}(ctx context.Context, data *model.{{.StructName}}) error {
	affect, err := repo.data.db.Context(ctx).Insert(data)
	if err != nil {
		return err
	}
	if affect != 1 {
		return errcode.DBUpdateNotAffected
	}
	return nil
}

func (repo *{{.StructName}}Repo) FindOne(ctx context.Context, id {{.PrimaryKeyType}}) (*model.{{.StructName}}, error) {
	var resp model.{{.StructName}}
	has, err := repo.data.db.Context(ctx).ID(id).Get(&resp)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errcode.DBDataNotFound
	}
	return &resp, nil
}

func (repo *{{.StructName}}Repo) Update(ctx context.Context, data *model.{{.StructName}}) error {
	_, err := repo.data.db.Context(ctx).ID(data.{{.PrimaryKeyName}}).Update(data)
	if err != nil {
		return err
	}

	return nil
}

func (repo *{{.StructName}}Repo) Delete(ctx context.Context, id {{.PrimaryKeyType}}) error {
	affect, err := repo.data.db.Context(ctx).ID(id).Delete(model.{{.StructName}}{})
	if err != nil {
		return err
	}
	if affect != 1 {
		return errcode.DBUpdateNotAffected
	}
	return nil
}

func (repo *{{.StructName}}Repo) List(ctx context.Context,in request.{{.StructName}}Req) (list []*model.{{.StructName}}, err error) {
	list = make([]*model.{{.StructName}}, 0)
	err = repo.data.db.Context(ctx).Limit(in.Limit()).Find(&list)
	return
}




`)

var tpl_service = fmt.Sprintf(`package service
import (
 	"context"
	"{{.GoMod}}/model"
	"{{.GoMod}}/internal/repo"
)

type {{.StructName}}Service struct {
	{{.DownLatterStructName}}Repo *{{.StructName}}Repo
}

func New{{.StructName}}Service({{.DownLatterStructName}}Repo *{{.StructName}}Repo) *{{.StructName}}Service {
	service := {{.StructName}}Service{}
	service.{{.DownLatterStructName}}Repo = {{.StructName}}Repo
	return service
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
