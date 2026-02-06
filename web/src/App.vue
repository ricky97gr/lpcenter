<template>
  <a-layout style="min-height: 100vh">
    <a-layout-header style="background: #001529; padding: 0 50px">
      <div style="color: white; font-size: 20px; font-weight: bold">
        插件授权中心
      </div>
    </a-layout-header>
    <a-layout>
      <a-layout-sider width="200" style="background: #fff">
        <a-menu
          mode="inline"
          :default-selected-keys="['1']"
          style="height: 100%"
          @click="handleMenuClick"
        >
          <a-menu-item key="1">
            <router-link to="/plugins">插件列表</router-link>
          </a-menu-item>
          <a-menu-item key="2">
            <router-link to="/license-request">申请授权</router-link>
          </a-menu-item>
          <a-menu-item key="3">
            <router-link to="/license-management">授权管理</router-link>
          </a-menu-item>
        </a-menu>
      </a-layout-sider>
      <a-layout-content style="padding: 24px">
        <div style="background: #fff; padding: 24px; min-height: 380px">
          <router-view />
        </div>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const handleMenuClick = ({ key }) => {
  const routes = {
    '1': '/plugins',
    '2': '/license-request',
    '3': '/license-management'
  }
  router.push(routes[key])
}
</script>

<style scoped>
.ant-layout-header {
  display: flex;
  align-items: center;
}
</style>
