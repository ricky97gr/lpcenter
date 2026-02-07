<template>
  <div class="main-layout">
    <div class="sider-container" :class="{ collapsed: collapsed }">
      <div class="logo">
        <span v-if="!collapsed">插件授权中心</span>
        <span v-else>LP</span>
      </div>
      <div class="menu-container">
        <a-menu
          :selected-keys="[selectedKey]"
          :open-keys="openKeys"
          :inline-collapsed="collapsed"
          mode="inline"
          theme="light"
          @click="handleMenuClick"
          @openChange="handleOpenChange"
        >
          <a-menu-item key="/dashboard">
            <template #icon>
              <DashboardOutlined />
            </template>
            <span>概览</span>
          </a-menu-item>

          <a-menu-item key="/versions">
            <template #icon>
              <AppstoreOutlined />
            </template>
            <span>授权类型管理</span>
          </a-menu-item>

          <a-menu-item key="/products">
            <template #icon>
              <DatabaseOutlined />
            </template>
            <span>产品管理</span>
          </a-menu-item>

          <a-menu-item key="/plugins">
            <template #icon>
              <AppstoreOutlined />
            </template>
            <span>插件管理</span>
          </a-menu-item>

          <a-menu-item key="/license-request">
            <template #icon>
              <FormOutlined />
            </template>
            <span>申请授权</span>
          </a-menu-item>

          <a-menu-item key="/license-management">
            <template #icon>
              <SafetyOutlined />
            </template>
            <span>授权管理</span>
          </a-menu-item>

          <a-menu-item key="/user-management">
            <template #icon>
              <UserOutlined />
            </template>
            <span>用户管理</span>
          </a-menu-item>
        </a-menu>
      </div>
      <div class="sider-footer">
        <MenuFoldOutlined
          v-if="!collapsed"
          class="trigger"
          @click="toggleCollapsed"
        />
        <MenuUnfoldOutlined
          v-else
          class="trigger"
          @click="toggleCollapsed"
        />
      </div>
    </div>

    <div class="content-wrapper" :class="{ collapsed: collapsed }">
      <div class="header">
        <div style="flex: 1; display: flex; align-items: center;">
          <div v-if="currentPageTitle" class="page-title-container">
            <span class="page-title">{{ currentPageTitle }}</span>
          </div>
        </div>
        <div class="user-info">
          <a-dropdown>
            <a class="ant-dropdown-link" @click.prevent>
              <UserOutlined /> 管理员
              <DownOutlined />
            </a>
            <template #overlay>
              <a-menu @click="handleMenuClick">
                <a-menu-item key="profile">
                  <UserOutlined /> 个人中心
                </a-menu-item>
                <a-menu-item key="settings">
                  <SettingOutlined /> 个人设置
                </a-menu-item>
                <a-menu-divider />
                <a-menu-item key="logout">
                  <LogoutOutlined /> 退出登录
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </div>
      <div class="content-area">
        <div style="background: #fff; margin: 0; padding: 10px; height: 100%; border-radius: 8px;">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  AppstoreOutlined,
  FormOutlined,
  SafetyOutlined,
  DatabaseOutlined,
  SettingOutlined,
  DownOutlined,
  LogoutOutlined,
  UserOutlined,
  DashboardOutlined
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { authAPI } from '../api'

const router = useRouter()
const route = useRoute()

const collapsed = ref(false)
const selectedKey = computed(() => route.path)
const openKeys = ref([])

const currentPageTitle = computed(() => {
  const path = route.path
  if (path === '/dashboard') return '概览'
  if (path === '/versions') return '授权类型管理'
  if (path === '/products') return '产品管理'
  if (path === '/plugins') return '插件管理'
  if (path === '/license-request') return '申请授权'
  if (path === '/license-management') return '授权管理'
  if (path === '/user-management') return '用户管理'
  return ''
})

const toggleCollapsed = () => {
  collapsed.value = !collapsed.value
}

const handleMenuClick = ({ key }) => {
  if (key === 'logout') {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    message.success('退出登录成功')
    router.push('/login')
  } else if (key === 'profile' || key === 'settings') {
    message.info('功能开发中')
  } else {
    router.push(key)
  }
}

const handleOpenChange = (keys) => {
  openKeys.value = keys
}
</script>

<style scoped>
.main-layout {
  height: 100vh;
  overflow: hidden;
  display: flex;
  width: 100vw;
  margin: 0;
  padding: 0;
}

.sider-container {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: 200px;
  z-index: 10;
  overflow: hidden;
  transition: width 0.3s;
  background: #fff;
  border-right: 1px solid #f0f0f0;
}

.sider-container.collapsed {
  width: 80px;
}

.content-wrapper {
  flex: 1;
  margin-left: 200px;
  transition: margin-left 0.3s;
  height: 100vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.content-wrapper.collapsed {
  margin-left: 80px;
}

.header {
  padding: 0 16px;
  background: #fff;
  display: flex;
  align-items: center;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  height: 60px;
  line-height: 60px;
  flex-shrink: 0;
}

.content-area {
  flex: 1;
  margin: 12px;
  padding: 0;
  background: #f0f2f5;
  overflow-y: auto;
  overflow-x: hidden;
  border-radius: 8px;
}

.content-area::-webkit-scrollbar {
  width: 8px;
}

.content-area::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.content-area::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

.content-area::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

.logo {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: bold;
  color: #1890ff;
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
}

.trigger {
  font-size: 14px;
  padding: 0 16px;
  cursor: pointer;
  transition: color 0.3s;
}

.trigger:hover {
  color: #1890ff;
}

.menu-container {
  height: calc(100vh - 112px);
  overflow-y: auto;
  overflow-x: hidden;
}

.sider-footer {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding: 12px 0;
  background: #fff;
  border-top: 1px solid #f0f0f0;
}

.menu-container::-webkit-scrollbar {
  width: 6px;
}

.menu-container::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.menu-container::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.menu-container::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

.user-info {
  display: flex;
  align-items: center;
}

.page-title-container {
  margin-right: 24px;
}

.page-title {
  font-size: 16px;
  font-weight: bold;
  color: #333;
}

:deep(.ant-menu-item),
:deep(.ant-menu-submenu-title) {
  font-size: 14px;
  padding-left: 16px !important;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

:deep(.ant-menu-submenu-title > span) {
  flex: 1;
  text-align: left;
}

:deep(.ant-menu-submenu-title > .anticon-down) {
  margin-left: auto;
}

:deep(.ant-menu-submenu > .ant-menu) {
  padding-left: 16px !important;
}

:deep(.ant-menu-submenu > .ant-menu > .ant-menu-item) {
  padding-left: 32px !important;
  margin-left: 0 !important;
}

:deep(.ant-menu-item-group > .ant-menu-item) {
  padding-left: 32px !important;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
