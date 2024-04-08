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

var tpl_dao = fmt.Sprintf(`package dao

import (
 	"context"
	"{{.Module}}/model"
)

type {{.StructName}}Repo struct {
	data *Data
}

func New{{.StructName}}Repo(data *Data) *{{.StructName}}Repo {
	repo := {{.StructName}}Repo{}
	repo.data = data
	return &repo
}

// Insert 创建{{.StructName}}
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
	_, err := repo.data.db.Context(ctx).ID(data.{{.PkName}}).Update(data)
	if err != nil {
		return err
	}

	return nil
}

func (repo *{{.StructName}}Repo) Delete(ctx context.Context, id {{.PkType}}) error {
	affect, err := repo.data.db.Context(ctx).ID(id).Delete(model.{{.StructName}}{})
	if err != nil {
		return err
	}
	if affect != 1 {
		return errcode.DBUpdateNotAffected
	}
	return nil
}

func (repo *{{.StructName}}Repo) List(ctx context.Context,in request.{{.StructName}}Search) (list []*model.{{.StructName}}, err error) {
	list = make([]*model.{{.StructName}}, 0)
	err = repo.data.db.Context(ctx).Limit(in.Limit()).Find(&list)
	return
}




`)