package genCode

import (
	"bytes"
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
