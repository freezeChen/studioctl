import axiosInstance from "@/api";



export async function loadGoInfo(target:string): Promise<string> {
    return axiosInstance.get("setting/loadGoInfo", {
        params: {target: target},
    });
}