package genCode

import (
	"bytes"
	"fmt"
	"github.com/freezeChen/studioctl/core/util"
	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli"
	template2 "html/template"
	"path"
	"text/template"
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
		tableMapper.RestPackage = path.Join("interface", "http/web", tableMapper.RouterPath, tableMapper.Module)

		tableMapper.LastModelPackage = tableMapper.Module + "Model"
		tableMapper.LastRequestPackage = tableMapper.Module + "Req"
		tableMapper.LastDaoPackage = tableMapper.Module + "Dao"
		tableMapper.LastServicePackage = tableMapper.Module + "Service"
	} else {
		tableMapper.ModelPackage = path.Join(tableMapper.GoMod, "model")
		tableMapper.ServicePackage = path.Join(tableMapper.GoMod, "service")
		tableMapper.DaoPackage = path.Join(tableMapper.GoMod, "dao")
		tableMapper.RequestPackage = path.Join(tableMapper.ModelPackage, "req")
		tableMapper.ResponsePackage = path.Join(tableMapper.ModelPackage, "res")
		tableMapper.RestPackage = path.Join(tableMapper.GoMod, "interface", "http/web", tableMapper.RouterPath)

		tableMapper.LastModelPackage = "model"
		tableMapper.LastRequestPackage = "req"
		tableMapper.LastDaoPackage = "dao"
		tableMapper.LastServicePackage = "service"
	}

	for _, field := range in.Fields {
		column := ColumnMapper{
			Name:       field.FieldJson,
			MapperName: field.FieldName,
			ZhName:     field.FieldZhName,
			JsonName:   util.FirstLower(util.PascalCase(field.FieldJson)),
			Type:       field.FieldType,
			Comment:    field.FieldComment,
			IsAuto:     field.IsAuto,
			SearchType: field.SearchType,
			IsKey:      field.IsKey,
			Show:       field.Show,
			Require:    field.Require,
			Tag:        template2.HTML(fmt.Sprintf("`json:\"%s\"`", field.FieldJson)),
		}

		if field.FieldType == "jsontime.JsonTime" {
			if field.SearchType != "" {
				tableMapper.ReqHasJsonTime = true
			}
			tableMapper.ModelHasJsonTime = true
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

	toString, _ := jsoniter.MarshalToString(in)
	println(toString)

	var out = new(CodeInfo)
	model, err := genModel(in)
	if err != nil {
		return nil, err
	}

	out.Codes = append(out.Codes, CodeItem{
		FileName: in.StructName + ".go",
		Path:     fmt.Sprintf("internal/%s", in.ModelPackage),
		Code:     model,
		Type:     1,
	})

	req, err := genReq(in)
	if err != nil {
		return nil, err
	}

	out.Codes = append(out.Codes, CodeItem{
		FileName: in.StructName + "Req.go",
		Path:     fmt.Sprintf("internal/%s", in.RequestPackage),
		Code:     req,
		Type:     1,
	})

	repo, err := genRepo(in)
	if err != nil {
		return nil, err
	}

	out.Codes = append(out.Codes, CodeItem{
		FileName: in.StructName + "Repo.go",
		Path:     fmt.Sprintf("internal/%s", in.DaoPackage),
		Code:     repo,
		Type:     1,
	})

	svc, err := genService(in)
	if err != nil {
		return nil, err
	}
	out.Codes = append(out.Codes, CodeItem{
		FileName: in.StructName + "Service.go",
		Path:     fmt.Sprintf("internal/%s", in.ServicePackage),
		Code:     svc,
		Type:     1,
	})

	rest, err := genRest(in)
	if err != nil {
		return nil, err
	}
	out.Codes = append(out.Codes, CodeItem{
		FileName: in.StructName + "Rest.go",
		Path:     fmt.Sprintf("internal/%s", in.RestPackage),
		Code:     rest,
		Type:     1,
	})

	webApi, err := genWebApi(in)
	if err != nil {
		return nil, err
	}
	out.Codes = append(out.Codes, CodeItem{
		FileName: in.DownLatterStructName + ".js",
		Path:     fmt.Sprintf(path.Join("src/api", in.Module)),
		Code:     webApi,
		Type:     2,
	})

	webForm, err := genWebForm(in)
	if err != nil {
		return nil, err
	}
	out.Codes = append(out.Codes, CodeItem{
		FileName: in.StructName + "form.vue",
		Path:     fmt.Sprintf(path.Join("src/view", in.Module, in.DownLatterStructName)),
		Code:     webForm,
		Type:     2,
	})
	web, err := genWeb(in)
	if err != nil {
		return nil, err
	}
	out.Codes = append(out.Codes, CodeItem{
		FileName: in.StructName + "web.vue",
		Path:     fmt.Sprintf(path.Join("src/view", in.Module, in.DownLatterStructName)),
		Code:     web,
		Type:     2,
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

func genWebApi(in TableMapper) (string, error) {
	parse, err := template.New("genRest").Parse(tpl_web_api)
	if err != nil {
		return "", err
	}
	var b = bytes.Buffer{}

	parse.Execute(&b, in)
	return b.String(), nil
}

func genWebForm(in TableMapper) (string, error) {
	parse, err := template.New("genRest").Parse(tpl_web_form)
	if err != nil {
		return "", err
	}
	var b = bytes.Buffer{}

	parse.Execute(&b, in)
	return b.String(), nil
}

func genWeb(in TableMapper) (string, error) {
	parse, err := template.New("genRest").Parse(tpl_web_view)
	if err != nil {
		return "", err
	}
	var b = bytes.Buffer{}

	parse.Execute(&b, in)
	return b.String(), nil
}
