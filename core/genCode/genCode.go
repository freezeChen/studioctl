package genCode

import (
	"bytes"
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

	NewQuery(url, dbType)

	NewServer(port)
}

func genModel(in PreviewReq) (string, error) {
	var tableMapper = TableMapper{
		GoMod:      "",
		TableName:  in.TableName,
		StructName: in.StructName,
		Comment:    in.Comment,
	}
	for _, field := range in.Fields {
		column := ColumnMapper{
			Name:       field.FieldJson,
			MapperName: field.FieldName,
			ZhName:     field.FieldZhName,
			Type:       field.FieldType,
			Comment:    field.FieldComment,
			IsAuto:     field.IsAuto,
			IsKey:      field.IsKey,
		}
		column.Tag = column.ColumnTag()

		tableMapper.Columns = append(tableMapper.Columns, column)
	}

	parse, _ := template.New("genModel").Parse(tpl_model)
	var b = bytes.Buffer{}

	parse.Execute(&b, tableMapper)
	return b.String(), nil
}
