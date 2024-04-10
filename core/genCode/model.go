package genCode

import (
	"github.com/freezeChen/studioctl/core/util"
	"html/template"
	"slices"
	"strings"
)

type PreviewReq struct {
	TableName  string         `json:"table_name,omitempty"`  //表名
	StructName string         `json:"struct_name,omitempty"` //实体名称
	FileName   string         `json:"file_name,omitempty"`   //文件名称
	Comment    string         `json:"comment"`               //表注释
	Module     string         `json:"module,omitempty"`      //模块(具体为创建一级目录)
	Fields     []PreviewField `json:"fields,omitempty"`
}

type PreviewField struct {
	FieldName    string `json:"field_name,omitempty"`    //字段名
	FieldZhName  string `json:"field_zh_name,omitempty"` //字段中文名
	FieldComment string `json:"field_comment,omitempty"` //字段备注
	FieldType    string `json:"field_type,omitempty"`    //字段类型
	FieldJson    string `json:"field_json,omitempty"`    //jsonTag
	Show         bool   `json:"show"`                    //是否显示
	Require      bool   `json:"require,omitempty"`       //是否必填(编辑)
	SearchType   string `json:"search_type,omitempty"`   //搜索条件(=,like,between)
	IsKey        bool   `json:"is_key"`                  //是否主键
	IsAuto       bool   `json:"is_auto"`                 //是否自增
}

type TableMapper struct {
	GoMod                string
	TableName            string
	StructName           string
	DownLatterStructName string
	TableZhName          string
	Comment              string
	PrimaryKeyType       string //主键类型
	PrimaryKeyName       string //主键名称
	Columns              []ColumnMapper
}

type ColumnMapper struct {
	Name       string
	MapperName string
	ZhName     string
	Type       string
	Comment    string
	IsKey      bool
	IsAuto     bool
	SearchType string
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

type CodeInfo struct {
	PrefixPath string     `json:"prefix_path"` //路径前缀
	PkgPackage string     `json:"pkg_package"` //自定义包路径
	Codes      []CodeItem `json:"codes"`
}

type CodeItem struct {
	FileName string `json:"file_name"` //文件名称
	Path     string `json:"path"`      //保存路径
	Code     string `json:"code"`      //代码
}
