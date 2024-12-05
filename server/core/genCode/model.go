package genCode

import (
	"fmt"
	"github.com/freezeChen/studioctl/core/util"
	"html/template"
	"slices"
	"strings"
)

type TableInfoRs struct {
	TableName  string         `json:"table_name,omitempty"`  //表名
	StructName string         `json:"struct_name,omitempty"` //实体名称
	FileName   string         `json:"file_name,omitempty"`   //文件名称
	Comment    string         `json:"comment"`               //表注释
	Module     string         `json:"module,omitempty"`      //模块(具体为创建一级目录)
	ChName     string         `json:"ch_name"`               //中文名称(用于备注提示)
	Fields     []PreviewField `json:"fields,omitempty"`
}

type PreviewReq struct {
	TableName  string `json:"table_name,omitempty"`  //表名
	StructName string `json:"struct_name,omitempty"` //实体名称
	FileName   string `json:"file_name,omitempty"`   //文件名称

	Comment       string         `json:"comment"`          //表注释
	Module        string         `json:"module,omitempty"` //模块(具体为创建一级目录)
	ChName        string         `json:"ch_name"`          //中文名称(用于备注提示)
	GoOutDir      string         `json:"go_out_dir"`       //go代码输出路径
	JsOutDir      string         `json:"js_out_dir"`       //js代码输出路径
	PackagePrefix string         `json:"package_prefix"`   //go包前缀
	RouterPath    string         `json:"router_path"`      //路由包路径
	Fields        []PreviewField `json:"fields,omitempty"`
}

type PreviewField struct {
	FieldName    string `json:"field_name,omitempty"`    //字段名
	FieldZhName  string `json:"field_zh_name,omitempty"` //字段中文名
	FieldComment string `json:"field_comment,omitempty"` //字段备注
	FieldType    string `json:"field_type,omitempty"`    //字段类型
	FieldJson    string `json:"field_json,omitempty"`    //jsonTag
	DictType     string `json:"dict_type,omitempty"`
	Show         bool   `json:"show"`                  //是否显示
	Require      bool   `json:"require,omitempty"`     //是否必填(编辑)
	SearchType   string `json:"search_type,omitempty"` //搜索条件(=,like,between)
	IsKey        bool   `json:"is_key"`                //是否主键
	IsAuto       bool   `json:"is_auto"`               //是否自增
}

type TableMapper struct {
	GoMod                string //go mod
	TableName            string //表面
	StructName           string //结构体名称
	DownLatterStructName string
	TableZhName          string //表中文名称
	Comment              string //表注释
	PrimaryKeyType       string //主键类型
	PrimaryKeyName       string //主键名称
	Module               string //模块名称
	RouterPath           string //路由包路径
	//用于方便代码生成
	ModelHasJsonTime   bool
	ReqHasJsonTime     bool
	ModelPackage       string
	LastModelPackage   string
	ServicePackage     string
	LastServicePackage string
	RestPackage        string
	DaoPackage         string
	LastDaoPackage     string
	RequestPackage     string
	LastRequestPackage string
	ResponsePackage    string

	Columns []ColumnMapper
}

type ColumnMapper struct {
	Name       string
	MapperName string
	ZhName     string
	JsonName   string
	Type       string
	DictType   string
	Comment    string
	Show       bool
	Require    bool
	IsKey      bool
	IsAuto     bool
	SearchType string
	Tag        template.HTML
}

func (c ColumnMapper) ColumnTag() template.HTML {
	var res []string
	tag := util.GormTag{}
	if c.IsKey {
		res = append(res, tag.PrimaryKey())
	}
	if c.IsAuto {
		res = append(res, tag.Auto())
	}
	if c.Name == "created_at" {
		res = append(res, tag.CreateTime())
	}
	if c.Name == "updated_at" {
		res = append(res, tag.UpdateTime())
	}

	slices.Sort(res)
	if len(res) == 0 {
		return template.HTML("")
	}
	return template.HTML(fmt.Sprintf(`%s:"%s"`, tag.TagName(), strings.Join(res, tag.Split())))
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
	Type     int32  `json:"type"`      //1:go 2:vue
}

type DictInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
