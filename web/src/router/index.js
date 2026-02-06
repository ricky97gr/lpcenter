import { createRouter, createWebHistory } from 'vue-router'
import PluginList from '../views/PluginList.vue'
import PluginDetail from '../views/PluginDetail.vue'
import LicenseRequest from '../views/LicenseRequest.vue'
import LicenseManagement from '../views/LicenseManagement.vue'

const routes = [
  {
    path: '/',
    redirect: '/plugins'
  },
  {
    path: '/plugins',
    name: 'PluginList',
    component: PluginList
  },
  {
    path: '/plugins/:id',
    name: 'PluginDetail',
    component: PluginDetail
  },
  {
    path: '/license-request',
    name: 'LicenseRequest',
    component: LicenseRequest
  },
  {
    path: '/license-management',
    name: 'LicenseManagement',
    component: LicenseManagement
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
