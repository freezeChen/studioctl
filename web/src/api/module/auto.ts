import axiosInstance from "..";

export interface TableRes {
    table_name: string;
    table_comment: string;
}



export interface Preview {
    table_name: string; //表名
    struct_name: string; //实体名称
    file_name: string; //文件名称
    module: string;//模块(具体为创建一级目录)
    comment: string; //表注释
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
    is_key: boolean;
    is_auto: boolean;
}


export async function getTables(): Promise<TableRes[]> {
    return axiosInstance.get("gen/tables");
}

export async function getTableColumns(
    params: string
): Promise<Preview> {
    return axiosInstance.get("gen/columns", {
        params: {table: params},
    });
}

export async function previewCode(params:Preview): Promise<string> {
    return axiosInstance.post("gen/preview",params)
}