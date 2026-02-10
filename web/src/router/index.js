import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '../layouts/MainLayout.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/upload',
    name: 'PluginUpload',
    component: () => import('../views/PluginUpload.vue'),
    meta: { title: '插件上传', public: true }
  },
  {
    path: '/',
    name: 'Layout',
    component: MainLayout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { title: '概览' }
      },
      {
        path: 'license-types',
        name: 'LicenseTypeManagement',
        component: () => import('../views/LicenseTypeManagement.vue'),
        meta: { title: '授权类型管理' }
      },
      {
        path: 'products',
        name: 'ProductManagement',
        component: () => import('../views/ProductManagement.vue'),
        meta: { title: '产品管理' }
      },
      {
        path: 'plugins',
        name: 'PluginManagement',
        component: () => import('../views/PluginManagement.vue'),
        meta: { title: '插件管理' }
      },
      {
        path: 'license-request',
        name: 'LicenseRequest',
        component: () => import('../views/LicenseRequest.vue'),
        meta: { title: '申请授权' }
      },
      {
        path: 'license-management',
        name: 'LicenseManagement',
        component: () => import('../views/LicenseManagement.vue'),
        meta: { title: '授权管理' }
      },
      {
        path: 'user-management',
        name: 'UserManagement',
        component: () => import('../views/UserManagement.vue'),
        meta: { title: '用户管理' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  if (to.meta.title) {
    document.title = `${to.meta.title} - 插件授权中心`
  }
  
  if (to.meta.public) {
    next()
  } else if (to.path !== '/login' && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router