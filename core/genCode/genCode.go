package genCode

import (
	"bytes"
	"fmt"
	"github.com/freezeChen/studioctl/core/util"
	"github.com/urfave/cli"
	"html/template"
)

const (
	FlagUrl  = "url"
	FlagPort = "port"
	FlagType = "dbType"
)

func CodeHandler(ctx *cli.Context) {
	url := ctx.String(FlagUrl)
	port := ctx.String(FlagPort)
	dbType := ctx.String(FlagType)

	moduleName, err := util.GetGoModuleName()
	if err != nil {
		panic(err)
	}
	c.GoMod = moduleName
	fmt.Println(moduleName)

	NewQuery(url, dbType)

	NewServer(port)
}

func parseTableColumns(table string, columns []Column) PreviewReq {
	var out PreviewReq
	out.TableName = table
	out.StructName = util.PascalCase(table)

	for _, column := range columns {
		out.Fields = append(out.Fields, PreviewField{
			FieldName:    util.PascalCase(column.ColumnName),
			FieldZhName:  column.ColumnComment,
			FieldComment: column.ColumnComment,
			FieldType:    util.SQLTypeToStructType(column.DataType),
			FieldJson:    column.ColumnName,
			Require:      false,
			Show:         true,
			SearchType:   "",
			IsKey:        column.IsKey,
			IsAuto:       column.IsAuto})
	}

	return out
}

func getTableMapper(in PreviewReq) TableMapper {
	var tableMapper = TableMapper{
		GoMod:                c.GoMod,
		TableName:            in.TableName,
		StructName:           in.StructName,
		DownLatterStructName: util.FirstLower(in.StructName),
		Comment:              in.Comment,
	}
	for _, field := range in.Fields {
		column := ColumnMapper{
			Name:       field.FieldJson,
			MapperName: field.FieldName,
			ZhName:     field.FieldZhName,
			Type:       field.FieldType,
			Comment:    field.FieldComment,
			IsAuto:     field.IsAuto,
			SearchType: field.SearchType,
			IsKey:      field.IsKey,
			Tag:        template.HTML(fmt.Sprintf("`json:\"%s\"`", field.FieldJson)),
		}
		column.Tag = column.ColumnTag()
		if field.IsKey {
			tableMapper.PrimaryKeyType = field.FieldType
		}

		tableMapper.Columns = append(tableMapper.Columns, column)
	}

	return tableMapper
}

func genCode(in TableMapper) (map[string]string, error) {
	model, err := genModel(in)
	if err != nil {
		return nil, err
	}
	req, err := genReq(in)
	if err != nil {
		return nil, err
	}

	repo, err := genRepo(in)
	if err != nil {
		return nil, err
	}
	svc, err := genService(in)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		in.StructName + ".go":        model,
		in.StructName + "Repo.go":    repo,
		in.StructName + "Req.go":     req,
		in.StructName + "Service.go": svc,
	}, nil
}

func genModel(in TableMapper) (string, error) {
	parse, _ := template.New("genModel").Parse(tpl_model)
	var b = bytes.Buffer{}
	parse.Execute(&b, in)
	return b.String(), nil
}

func genReq(in TableMapper) (string, error) {
	parse, _ := template.New("genreq").Parse(tpl_modelReq)
	var b = bytes.Buffer{}
	parse.Execute(&b, in)
	return b.String(), nil
}

func genRepo(in TableMapper) (string, error) {
	parse, _ := template.New("genRepo").Parse(tpl_Repo)
	var b = bytes.Buffer{}
	parse.Execute(&b, in)
	return b.String(), nil
}

func genService(in TableMapper) (string, error) {
	parse, err := template.New("genSerivce").Parse(tpl_service)
	if err != nil {
		return "", err
	}
	var b = bytes.Buffer{}

	parse.Execute(&b, in)
	return b.String(), nil
}
