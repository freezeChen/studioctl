import axios from "axios";

interface HttpResponse<T = unknown> {
    code: number;
    msg: string;
    data: T;
}

const axiosInstance = axios.create({
    baseURL: "http://localhost:8888/",
});

axiosInstance.interceptors.request.use((req) => {
    return req;
});

axiosInstance.interceptors.response.use((response) => {
    console.log(response.data);
    const data = response.data;
    if (data.code != 0) {
        return Promise.reject(Error(data.msg));
    }
    return data.data;
});

export default axiosInstance;
