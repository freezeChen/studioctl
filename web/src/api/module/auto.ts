import axiosInstance from "..";

export interface TableRes {
    table_name: string;
    table_comment: string;
}


export interface GenTableInfoRes {
    table_name: string; //表名
    struct_name: string; //实体名称
    file_name: string; //文件名称
    module: string;//模块(具体为创建一级目录)
    comment: string; //表注释
    ch_name: string;
    fields: Array<PreviewField>;
}

export interface Preview {
    table_name: string; //表名
    struct_name: string; //实体名称
    file_name: string; //文件名称
    module: string;//模块(具体为创建一级目录)
    comment: string; //表注释
    ch_name: string;
    go_out_dir: string; //go代码输出路径
    js_out_dir: string; //前端代码输出路径
    package_prefix: string; //go 包前缀
    router_path: string;
    fields: Array<PreviewField>;
}

export interface PreviewField {
    field_name: string; //字段名称
    field_zh_name: string; //字段中文名
    field_comment: string; //字段注释
    field_type: string; //字段类型
    field_json: string; //字段json(sql字段)
    require: boolean; //是否必须
    search_type: string; ////搜索条件(=,like,between)
    dict_type: string;
    is_key: boolean;
    is_auto: boolean;
}

export interface PreviewRes {
    prefix_path: string; //路径前缀
    pkg_package: string;
    codes: Array<CodeItem>;
}

export interface CodeItem {
    file_name: string;
    path: string;
    code: string;
}

export interface DictInfo {
    name: string;
    type: string;
}


export async function getTables(): Promise<TableRes[]> {
    return axiosInstance.get("gen/tables");
}

export async function getTableColumns(
    params: string
): Promise<GenTableInfoRes> {
    return axiosInstance.get("gen/columns", {
        params: {table: params},
    });
}

export async function previewCode(params: Preview): Promise<PreviewRes> {
    return axiosInstance.post("gen/preview", params)
}

export async function downloadCode(params: Preview): Promise<any> {
    return axiosInstance.post("gen/download", params)
}

export async function getDictList(): Promise<DictInfo[]>{
    return axiosInstance.get("getDictList")
}