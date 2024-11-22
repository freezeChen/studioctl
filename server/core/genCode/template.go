package genCode

import "fmt"

var tpl_model = fmt.Sprintf(`package {{.LastModelPackage}}

{{- if .ModelHasJsonTime}}

import "gitee.com/zx/ace-library/lib/jsontime"
{{- end}}
// {{.StructName}} {{.Comment}}
type {{.StructName}} struct{
{{range .Columns}}		{{.MapperName}}	{{.Type}}	%s{{.Tag}} json:"{{.Name}}"{{if gt (len .ZhName) 0}} comment:"{{.ZhName}}"{{end}}%s {{if gt (len .Comment) 0}}// {{.Comment}}{{end}}
{{end}}	}


func ({{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
`, "`", "`")

var tpl_modelReq = fmt.Sprintf(`package {{if .Module}}{{.Module}}Req{{else}}req{{end}}

import (
	"{{.GoMod}}/model/common/request"
)

// {{.StructName}} {{.Comment}}
type {{.StructName}}ListReq struct{
{{range .Columns}}{{if eq .SearchType "between"}}	Start{{.MapperName}}	{{- if eq .Type "jsontime.JsonTime"}}	string{{- else}}	{{.Type}}{{- end}}	%sjson:"start_{{.Name}}"%s{{printf "\n"}}	End{{.MapperName}}	{{- if eq .Type "jsontime.JsonTime"}}	string{{- else}}	{{.Type}}{{- end}}	%sjson:"end_{{.Name}}"%s{{printf "\n"}}{{end}}{{if eq .SearchType "like" "="}}	{{.MapperName}}	{{- if eq .Type "jsontime.JsonTime"}}	string{{- else}}	{{.Type}}{{- end}} %sjson:"{{.Name}}"%s{{printf "\n"}}{{end}}{{end}}	request.PageInfo
}



`, "`", "`", "`", "`", "`", "`")

var tpl_Repo = fmt.Sprintf(`package {{if .Module}}{{.Module}}Dao{{else}}dao{{end}}

import (
	{{if .Module}}"{{.GoMod}}/dao"{{end}}
	"{{.GoMod}}/{{.ModelPackage}}"
	"{{.GoMod}}/{{.RequestPackage}}"
	"context"
	"errors"
	"gorm.io/gorm"
)

type {{.StructName}}Repo struct {
	data *{{if .Module}}dao.{{end}}Data
}

func New{{.StructName}}Repo(data *{{if .Module}}dao.{{end}}Data) *{{.StructName}}Repo {
	repo := {{.StructName}}Repo{}
	repo.data = data
	return &repo
}

// Create{{.StructName}} 创建{{.TableZhName}}
func (repo *{{.StructName}}Repo) Create{{.StructName}}(ctx context.Context, data *{{.LastModelPackage}}.{{.StructName}}) error {
	err := repo.data.GetDB().WithContext(ctx).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *{{.StructName}}Repo) Get{{.StructName}}(ctx context.Context, id {{.PrimaryKeyType}}) (*{{.LastModelPackage}}.{{.StructName}}, error) {
	var resp {{.LastModelPackage}}.{{.StructName}}
	err := repo.data.GetDB().WithContext(ctx).First(&resp, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errcode.DBDataNotFound
	}
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (repo *{{.StructName}}Repo) Update{{.StructName}}(ctx context.Context, data *{{.LastModelPackage}}.{{.StructName}}) error {
	return repo.data.GetDB().WithContext(ctx).Where("id = ?",data.{{.PrimaryKeyName}}).Updates(data).Error
}

func (repo *{{.StructName}}Repo) Delete{{.StructName}}(ctx context.Context, id {{.PrimaryKeyType}}) error {
	return repo.data.GetDB().WithContext(ctx).Delete({{.LastModelPackage}}.{{.StructName}}{},id).Error
}

func (repo *{{.StructName}}Repo) Delete{{.StructName}}ByIds(ctx context.Context, ids []{{.PrimaryKeyType}}) error {
	return repo.data.GetDB().WithContext(ctx).Delete({{.LastModelPackage}}.{{.StructName}}{},ids).Error
}


func (repo *{{.StructName}}Repo) Get{{.StructName}}List(ctx context.Context,in {{.LastRequestPackage}}.{{.StructName}}ListReq) (list []*{{.LastModelPackage}}.{{.StructName}},count int64, err error) {
	list = make([]*{{.LastModelPackage}}.{{.StructName}}, 0)
	tx:=repo.data.GetDB().WithContext(ctx)
	{{- range .Columns}}{{- if eq .SearchType "between"}}
	if in.Start{{.MapperName}} != "" {
		tx = tx.Where("{{.Name}} > ?",in.Start{{.MapperName}})
	}
	if in.End{{.MapperName}} != "" {
		tx = tx.Where("{{.Name}} < ?",in.Start{{.MapperName}})
	}
{{- end}}{{- if eq .SearchType "like"}}
	if in.{{.MapperName}} != "" {
		tx = tx.Where("{{.Name}} like ?","%%"+in.{{.MapperName}}+"%%")
	}
{{- end}}{{- if eq .SearchType "="}}
	if in.{{.MapperName}} != {{- if eq .Type "string"}}""{{- else}}0{{- end}} {
		tx = tx.Where("{{.Name}} = ?",in.{{.MapperName}})
	}
{{- end}}{{- end}}

	tx= tx.Limit(in.Limit()).Offset(in.Offset())
	err = tx.Find(&list).Error
	if err != nil {
		return
	}
	err = tx.Count(&count).Error
	return
}
`)

var tpl_service = fmt.Sprintf(`package {{.LastServicePackage}}
import (
 	"context"
	"{{.GoMod}}/{{.DaoPackage}}"
	"{{.GoMod}}/model/common/response"
	"{{.GoMod}}/{{.ModelPackage}}"
	"{{.GoMod}}/{{.RequestPackage}}"
)

type {{.StructName}}Service struct {
	{{.DownLatterStructName}}Repo *{{.LastDaoPackage}}.{{.StructName}}Repo
}

func New{{.StructName}}Service({{.DownLatterStructName}}Repo *{{.LastDaoPackage}}.{{.StructName}}Repo) *{{.StructName}}Service {
	service := {{.StructName}}Service{}
	service.{{.DownLatterStructName}}Repo = {{.DownLatterStructName}}Repo
	return &service
}

// Create{{.StructName}} 创建{{.TableZhName}}
func (svc *{{.StructName}}Service) Create{{.StructName}}(ctx context.Context, data *{{.LastModelPackage}}.{{.StructName}}) error {
	return svc.{{.DownLatterStructName}}Repo.Create{{.StructName}}(ctx, data)
}

func (svc *{{.StructName}}Service) Get{{.StructName}}(ctx context.Context, id {{.PrimaryKeyType}}) (*{{.LastModelPackage}}.{{.StructName}}, error) {
	return svc.{{.DownLatterStructName}}Repo.Get{{.StructName}}(ctx, id)
}

func (svc *{{.StructName}}Service) Update{{.StructName}}(ctx context.Context, data *{{.LastModelPackage}}.{{.StructName}}) error {
	return svc.{{.DownLatterStructName}}Repo.Update{{.StructName}}(ctx, data)
}

func (svc *{{.StructName}}Service) Delete{{.StructName}}(ctx context.Context, id {{.PrimaryKeyType}}) error {
	return svc.{{.DownLatterStructName}}Repo.Delete{{.StructName}}(ctx, id)
}

func (svc *{{.StructName}}Service) Delete{{.StructName}}ByIds(ctx context.Context, ids []{{.PrimaryKeyType}}) error {
	return svc.{{.DownLatterStructName}}Repo.Delete{{.StructName}}ByIds(ctx, ids)
}


func (svc *{{.StructName}}Service) Get{{.StructName}}List(ctx context.Context,in {{.LastRequestPackage}}.{{.StructName}}ListReq) (*response.PageResult, error) {
	list,count,err:= svc.{{.DownLatterStructName}}Repo.Get{{.StructName}}List(ctx, in)
	if err != nil {
		return nil,err
	}
	return &response.PageResult{Total: count, List: list,Page: in.Page,PageSize: in.PageSize},nil
}
`)

var tpl_api = fmt.Sprintf(`package {{.Module}}
import (
	"{{.GoMod}}/{{.ModelPackage}}"
	"{{.GoMod}}/{{.RequestPackage}}"
	"{{.GoMod}}/{{.ServicePackage}}"
	"gitee.com/zx/ace-library/jsonresult"
	"github.com/gin-gonic/gin"
)

type {{.StructName}}Rest struct {
	{{.DownLatterStructName}}Service *{{.LastServicePackage}}.{{.StructName}}Service
	jsonresult.JsonResult
}

func New{{.StructName}}Rest({{.DownLatterStructName}}Service *{{.LastServicePackage}}.{{.StructName}}Service) *{{.StructName}}Rest {
	rest := {{.StructName}}Rest{}
	rest.{{.DownLatterStructName}}Service = {{.DownLatterStructName}}Service
	return &rest
}

func (rest {{.StructName}}Rest) Router(route *gin.RouterGroup) {
	group:=route.Group("{{.DownLatterStructName}}")
	group.POST("create{{.StructName}}",rest.create{{.StructName}})
	group.PUT("update{{.StructName}}",rest.update{{.StructName}})
	group.DELETE("delete{{.StructName}}",rest.delete{{.StructName}})
	group.DELETE("delete{{.StructName}}ByIds",rest.delete{{.StructName}}ByIds)
	group.POST("get{{.StructName}}List",rest.get{{.StructName}}List)
}

// create{{.StructName}} 创建{{.TableZhName}}
// @Tags {{.TableZhName}}
// @Summary 创建{{.TableZhName}}
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body {{.LastModelPackage}}.{{.StructName}} true "创建{{.TableZhName}}"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router	{{- if .Module}} /{{.Module}}{{- end}} /{{.DownLatterStructName}}/create{{.StructName}} [post]
func (rest {{.StructName}}Rest) create{{.StructName}}(ctx *gin.Context) {
	var param {{.LastModelPackage}}.{{.StructName}}
	if err := ctx.ShouldBind(&param); err != nil {
		rest.Error(ctx, err)
		return
	}
	err := rest.{{.DownLatterStructName}}Service.Create{{.StructName}}(ctx, &param)
	rest.Error(ctx, err)
}

// update{{.StructName}} 修改{{.TableZhName}}
// @Tags {{.TableZhName}}
// @Summary 修改{{.TableZhName}}
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body {{.LastModelPackage}}.{{.StructName}} true "修改{{.TableZhName}}"
// @Success 200 {object} response.Response{msg=string} "修改成功"
// @Router	{{- if .Module}} /{{.Module}}{{- end}} /{{.DownLatterStructName}}/update{{.StructName}} [put]
func (rest {{.StructName}}Rest) update{{.StructName}}(ctx *gin.Context) {
	var param {{.LastModelPackage}}.{{.StructName}}
	if err := ctx.ShouldBind(&param); err != nil {
		rest.Error(ctx, err)
		return
	}
	err := rest.{{.DownLatterStructName}}Service.Update{{.StructName}}(ctx, &param)
	rest.Error(ctx, err)
}

// delete{{.StructName}} 删除{{.TableZhName}}
// @Tags {{.TableZhName}}
// @Summary 删除{{.TableZhName}}
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id query int  true "id"
// @Success 200 {object} response.Response{msg=string} "修改成功"
// @Router{{- if .Module}} /{{.Module}}{{- end}} /{{.DownLatterStructName}}/update{{.StructName}} [put]
func (rest {{.StructName}}Rest) delete{{.StructName}}(ctx *gin.Context) {
	var param struct{
	{{.PrimaryKeyName}} {{.PrimaryKeyType}}` + "`" + `binding:"required"` + "`" + `
	}
	if err := ctx.ShouldBind(&param); err != nil {
		rest.Error(ctx, err)
		return
	}
	err := rest.{{.DownLatterStructName}}Service.Delete{{.StructName}}(ctx, param.{{.PrimaryKeyName}})
	rest.Error(ctx, err)
}

// delete{{.StructName}}ByIds 批量删除{{.TableZhName}}
// @Tags {{.TableZhName}}
// @Summary 批量删除{{.TableZhName}}
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ids query int  true "ids"
// @Success 200 {object} response.Response{msg=string} "修改成功"
// @Router	{{- if .Module}} /{{.Module}}{{- end}} /{{.DownLatterStructName}}/update{{.StructName}}ByIds [put]
func (rest {{.StructName}}Rest) delete{{.StructName}}ByIds(ctx *gin.Context) {
	var param struct{
	{{.PrimaryKeyName}}s []{{.PrimaryKeyType}} ` + "`" + `binding:"required"` + "`" + `
	}
	if err := ctx.ShouldBind(&param); err != nil {
		rest.Error(ctx, err)
		return
	}
	err := rest.{{.DownLatterStructName}}Service.Delete{{.StructName}}ByIds(ctx, param.{{.PrimaryKeyName}}s)
	rest.Error(ctx, err)
}

// get{{.StructName}}List 获取{{.TableZhName}}列表
// @Tags {{.TableZhName}}
// @Summary 获取{{.TableZhName}}列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body {{.LastRequestPackage}}.{{.StructName}}ListReq true "获取{{.TableZhName}}列表"
// @Success 200 {object} response.Response{data={{.LastModelPackage}}.{{.StructName}}} 
// @Router	{{- if .Module}} /{{.Module}}{{- end}} /{{.DownLatterStructName}}/get{{.StructName}}List [post]
func (rest {{.StructName}}Rest) get{{.StructName}}List(ctx *gin.Context) {
	var param {{.LastRequestPackage}}.{{.StructName}}ListReq
	if err := ctx.ShouldBind(&param); err != nil {
		rest.Error(ctx, err)
		return
	}
	out,err := rest.{{.DownLatterStructName}}Service.Get{{.StructName}}List(ctx, param)
	rest.Json(ctx,out, err)
}



`)

var tpl_web_api = `import service from '@/utils/request'

// 创建{{.TableZhName}}
export const create{{.StructName}}=(data)=>{
    return service({
     url: '{{- if .Module}}/{{.Module}}{{- end}}/{{.DownLatterStructName}}/create{{.StructName}}',
     method: 'post',
     data
    })
}


// 修改{{.TableZhName}}
export const update{{.StructName}}=(data)=>{
	return service({
	 url: '{{- if .Module}}/{{.Module}}{{- end}}/{{.DownLatterStructName}}/update{{.StructName}}',
	 method: 'put',
	 data
	})
}

// 删除{{.TableZhName}}
export const delete{{.StructName}}=(params)=>{
	return service({
	 url: '{{- if .Module}}/{{.Module}}{{- end}}/{{.DownLatterStructName}}/delete{{.StructName}}',
	 method: 'delete',
	 params
	})
}

// 批量删除{{.TableZhName}}
export const delete{{.StructName}}ByIds=(params)=>{
 	return service({
	 url: '{{- if .Module}}/{{.Module}}{{- end}}/{{.DownLatterStructName}}/delete{{.StructName}}ByIds',
	 method: 'delete',
	 params
	})
}


export const get{{.StructName}} = (params) => {
	return service({
	 url: '{{- if .Module}}/{{.Module}}{{- end}}/{{.DownLatterStructName}}/get{{.StructName}}',
	 method: 'post',
	 params
})
}

// 分页获取{{.TableZhName}}列表
export const get{{.StructName}}List = (params) => {
	return service({
	 url: '{{- if .Module}}/{{.Module}}{{- end}}/{{.DownLatterStructName}}/get{{.StructName}}List',
	 method: 'get',
	 params
})
}



`

var tpl_web_form = `<template>
    <div>
        <div class="gva-form-box">
            <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
                {{- range .Columns}}
                    <el-form-item label="{{.ZhName}}:" prop="{{.JsonName}}">
                        {{- if .Show}}
                            {{- if eq .Type "bool"}}
                                <el-switch v-model="formData.{{.JsonName}}" active-color="#13ce66"
                                           inactive-color="#ff4949"
                                           active-text="是" inactive-text="否" clearable></el-switch>
                            {{- end}}
                            {{- if eq .Type "string"}}
                                <el-input v-model="formData.{{.JsonName}}"
                                          placeholder="请输入{{.ZhName}}"></el-input>
                            {{- end}}
                            {{- if eq .Type "int" }}
                                <el-input v-model.number="formData.{{ .JsonName }}"
                                          placeholder="请输入{{.ZhName}}"/>
                            {{- end }}
                            {{- if eq .Type "jsontime.JsonTime" }}
                     <el-date-picker v-model="formData.{{ .JsonName }}" type="date" placeholder="选择日期""></el-date-picker>
                            {{- end }}
                        {{- end}}
                    </el-form-item>
                {{- end}}
                <el-form-item>
                    <el-button type="primary" @click="save">保存</el-button>
                    <el-button type="primary" @click="back">返回</el-button>
                </el-form-item>
            </el-form>
        </div>

    </div>
</template>

<script setup>
    import {
        create{{.StructName}},
        update{{.StructName}},
        find{{.StructName}}
    } from '@/plugin/{{.Module}}/api/{{.StructName}}'

    defineOptions({
        name: '{{.StructName}}Form'
    })

    // 自动获取字典
    import {getDictFunc} from '@/utils/format'
    import {useRoute, useRouter} from "vue-router"
    import {ElMessage} from 'element-plus'
    import {ref, reactive} from 'vue'

    const route = useRoute()
    const router = useRouter()

    const type = ref('')

    const formData = ref({
{{- range  .Columns}}
   
        {{- if eq .Type "bool"}}
        {{.JsonName}}: false,
        {{- end}}
        {{- if eq .Type "string"}}
        {{.JsonName}}: '',
        {{- end}}
        {{- if eq .Type "int32"}}
        {{.JsonName}}: 0,
        {{- end}}
        {{- if eq .Type "int64"}}
        {{.JsonName}}: 0,
        {{- end}}
        {{- if eq .Type "float32"}}
        {{.JsonName}}: 0,
        {{- end}}
        {{- if eq .Type "jsontime.JsonTime"}}
        {{.JsonName}}: new Date(),
        {{- end}}

{{- end}}
    })
    
// 验证规则
const rule = reactive({
    {{- range .Columns}}
    {{- if eq .Require true}}
    {{.JsonName}} : [{
        required: true,
        message: '不能为空',
        trigger: ['input','blur'],
    }],
    {{- end}}
      {{- end}}
})
    
const elFormRef = ref()

    // 初始化方法
const init = async () => {
    // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
        const res = await get{{.StructName}}({ ID: route.query.id })
        if (res.code === 0) {
            formData.value = res.data
            type.value = 'update'
        }
    } else {
        type.value = 'create'
    }

}

init()

const save = async () => {
        elFormRef.value?.validate( async (valid) => {
            if (!valid) return
            let res
            switch (type.value) {
                case 'create':
                    res = await create{{.StructName}}(formData.value)
                    break
                case 'update':
                    res = await update{{.StructName}}(formData.value)
                    break
                default:
                    res = await create{{.StructName}}(formData.value)
                    break
            }
            if (res.code === 0){
                ElMessage({
                    message: '创建/保存成功',
                    type: 'success',
                })

        }
        })
}

const back = () => {
    router.go(-1)
}

    


</script>

<style >
</style>`

var tpl_web_view = `<template>
    <div>
        <div class="gva-search-box">
            <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
                {{- range .Columns}}
                <el-form-item label="{{.ZhName}}" prop="{{.JsonName}}">
                    {{- if eq .Type "float64" "int"}}
                        {{if eq .SearchType "between"}}
                            
                            {{end}}
                        {{- end}}
                {{- end}}
            </el-form>
        </div>
    </div>
</template>
`
