import axiosInstance from "..";

export interface TableRes {
    table_name: string;
    table_comment: string;
}

export interface TableColumnRes {
    columnComment: string;
    columnName: string;
    dataType: string;
    dataTypeLong: string;
    isAuto: boolean;
    isKey: boolean;
    search: string;
    show: boolean;
}

export async function getTables(): Promise<TableRes[]> {
    return axiosInstance.get("getTables");
}

export async function getTableColumns(
    params: string
): Promise<TableColumnRes[]> {
    return axiosInstance.get("getTableColumns", {
        params: {tableName: params},
    });
}
