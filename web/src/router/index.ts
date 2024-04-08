import {createRouter, createWebHashHistory, RouteRecordRaw} from "vue-router";

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        component: () => import('@/pages/home/index.vue'),
        children: [
            {
                path: '/gencode',
                component: () => import('@/pages/gencode/index.vue')
            },
            {path: '/setting', component: () => import('@/pages/setting/index.vue')}
        ]
    }
]


const router = createRouter({
    history: createWebHashHistory(),
    routes,
});

export default router
