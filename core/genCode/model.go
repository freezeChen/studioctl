package genCode

import (
	"github.com/freezeChen/studioctl/core/util"
	"html/template"
	"slices"
	"strings"
)

type PreviewReq struct {
	TableName  string `json:"table_name,omitempty"`  //表名
	StructName string `json:"struct_name,omitempty"` //实体名称
	FileName   string `json:"file_name,omitempty"`   //文件名称
	Comment    string `json:"comment"`               //表注释
	Module     string `json:"module,omitempty"`      //模块(具体为创建一级目录)
	Fields     []struct {
		FieldName    string `json:"field_name,omitempty"`    //字段名
		FieldZhName  string `json:"field_zh_name,omitempty"` //字段中文名
		FieldComment string `json:"field_comment,omitempty"` //字段备注
		FieldType    string `json:"field_type,omitempty"`    //字段类型
		FieldJson    string `json:"field_json,omitempty"`    //jsonTag
		Require      bool   `json:"require,omitempty"`       //是否必填(编辑)
		SearchType   string `json:"search_type,omitempty"`   //搜索条件(=,like,between)
		IsKey        bool   `json:"is_key"`                  //是否主键
		IsAuto       bool   `json:"is_auto"`                 //是否自增
	} `json:"fields,omitempty"`
}

type TableMapper struct {
	GoMod      string
	TableName  string
	StructName string
	Comment    string
	Columns    []ColumnMapper
}

type ColumnMapper struct {
	Name       string
	MapperName string
	ZhName     string
	Type       string
	Comment    string
	IsKey      bool
	IsAuto     bool
	Tag        template.HTML
}

func (c ColumnMapper) ColumnTag() template.HTML {
	var res []string
	tag := util.NewTag()
	if c.IsKey {
		res = append(res, tag.PrimaryKey())
	}
	if c.IsAuto {
		res = append(res, tag.Auto())
	}

	slices.Sort(res)
	return template.HTML(strings.Join(res, tag.Split()))
}

type Table struct {
	TableName    string `json:"table_name"`
	TableComment string `json:"table_comment"`
}

type Column struct {
	DataType      string `json:"data_type"`
	ColumnName    string `json:"column_name" `
	DataTypeLong  string `json:"data_type_long" `
	ColumnComment string `json:"column_comment"`
	IsKey         bool   `json:"is_key"`  //是否主键
	IsAuto        bool   `json:"is_auto"` //是否自增
}
