package genCode

import (
	"bytes"
	"fmt"
	"github.com/freezeChen/studioctl/core/util"
	"github.com/urfave/cli"
	"html/template"
	"path"
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

func parseTableColumns(table string, columns []Column) (TableInfoRs, error) {
	var out TableInfoRs
	out.TableName = table
	out.FileName = table
	out.StructName = util.PascalCase(table)

	tableInfo, err := query.GetTableInfo(table)
	if err != nil {
		return TableInfoRs{}, err
	}
	out.ChName = tableInfo.TableComment
	if out.ChName == "" {
		out.ChName = tableInfo.TableName
	}

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

	return out, nil
}

func getTableMapper(in PreviewReq) TableMapper {
	var tableMapper = TableMapper{
		GoMod:                in.PackagePrefix,
		TableName:            in.TableName,
		StructName:           in.StructName,
		TableZhName:          in.ChName,
		DownLatterStructName: util.FirstLower(in.StructName),
		Comment:              in.Comment,
		Module:               in.Module,
	}
	if tableMapper.Module != "" {
		tableMapper.ModelPackage = path.Join("model", tableMapper.Module+"Model")
		tableMapper.ServicePackage = path.Join("service", tableMapper.Module+"Service")
		tableMapper.DaoPackage = path.Join("dao", tableMapper.Module+"Dao")
		tableMapper.RequestPackage = path.Join(tableMapper.ModelPackage, tableMapper.Module+"Req")

		tableMapper.ResponsePackage = path.Join(tableMapper.ModelPackage, tableMapper.Module+"Res")
		tableMapper.RestPackage = path.Join(tableMapper.RouterPath, tableMapper.Module)

		tableMapper.LastModelPackage = tableMapper.Module + "Model"
		tableMapper.LastRequestPackage = tableMapper.Module + "Req"

	} else {
		tableMapper.ModelPackage = path.Join(tableMapper.GoMod, "model")
		tableMapper.ServicePackage = path.Join(tableMapper.GoMod, "service")
		tableMapper.DaoPackage = path.Join(tableMapper.GoMod, "dao")
		tableMapper.RequestPackage = path.Join(tableMapper.ModelPackage, "req")
		tableMapper.ResponsePackage = path.Join(tableMapper.ModelPackage, "res")
		tableMapper.RestPackage = path.Join(tableMapper.GoMod, tableMapper.RouterPath)

		tableMapper.LastModelPackage = "model"
		tableMapper.LastRequestPackage = "req"
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
			tableMapper.PrimaryKeyName = field.FieldName
		}

		tableMapper.Columns = append(tableMapper.Columns, column)
	}

	return tableMapper
}

func genCode(in TableMapper) (*CodeInfo, error) {
	var out = new(CodeInfo)
	model, err := genModel(in)
	if err != nil {
		return nil, err
	}

	out.Codes = append(out.Codes, CodeItem{
		FileName: in.StructName + ".go",
		Path:     fmt.Sprintf("internal/%s", in.ModelPackage),
		Code:     model,
	})

	req, err := genReq(in)
	if err != nil {
		return nil, err
	}

	out.Codes = append(out.Codes, CodeItem{
		FileName: in.StructName + "Req.go",
		Path:     fmt.Sprintf("internal/%s", in.RequestPackage),
		Code:     req,
	})

	repo, err := genRepo(in)
	if err != nil {
		return nil, err
	}

	out.Codes = append(out.Codes, CodeItem{
		FileName: in.StructName + "Repo.go",
		Path:     fmt.Sprintf("internal/%s", in.DaoPackage),
		Code:     repo,
	})

	svc, err := genService(in)
	if err != nil {
		return nil, err
	}
	out.Codes = append(out.Codes, CodeItem{
		FileName: in.StructName + "Service.go",
		Path:     fmt.Sprintf("internal/%s", in.ServicePackage),
		Code:     svc,
	})

	rest, err := genRest(in)
	if err != nil {
		return nil, err
	}
	out.Codes = append(out.Codes, CodeItem{
		FileName: in.StructName + "Rest.go",
		Path:     "http",
		Code:     rest,
	})

	return out, nil
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

func genRest(in TableMapper) (string, error) {
	parse, err := template.New("genRest").Parse(tpl_api)
	if err != nil {
		return "", err
	}
	var b = bytes.Buffer{}

	parse.Execute(&b, in)
	return b.String(), nil
}
