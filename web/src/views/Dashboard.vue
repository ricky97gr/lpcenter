<template>
  <div class="dashboard">
    <a-row :gutter="[16, 16]">
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card :loading="loading">
          <a-statistic
            title="产品总数"
            :value="stats.totalProducts"
            :value-style="{ color: '#3f8600' }"
          >
            <template #prefix>
              <DatabaseOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card :loading="loading">
          <a-statistic
            title="版本总数"
            :value="stats.totalVersions"
            :value-style="{ color: '#1890ff' }"
          >
            <template #prefix>
              <AppstoreOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card :loading="loading">
          <a-statistic
            title="插件总数"
            :value="stats.totalPlugins"
            :value-style="{ color: '#cf1322' }"
          >
            <template #prefix>
              <AppstoreOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card :loading="loading">
          <a-statistic
            title="用户总数"
            :value="stats.totalUsers"
            :value-style="{ color: '#722ed1' }"
          >
            <template #prefix>
              <UserOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="[16, 16]" style="margin-top: 16px">
      <a-col :xs="24" :sm="12" :lg="8">
        <a-card title="授权统计" :loading="loading">
          <a-row :gutter="16">
            <a-col :span="8">
              <a-statistic
                title="待审核"
                :value="stats.pendingLicenses"
                :value-style="{ color: '#faad14' }"
              />
            </a-col>
            <a-col :span="8">
              <a-statistic
                title="已批准"
                :value="stats.approvedLicenses"
                :value-style="{ color: '#52c41a' }"
              />
            </a-col>
            <a-col :span="8">
              <a-statistic
                title="已拒绝"
                :value="stats.rejectedLicenses"
                :value-style="{ color: '#f5222d' }"
              />
            </a-col>
          </a-row>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="8">
        <a-card title="下载统计" :loading="loading">
          <a-statistic
            title="总下载量"
            :value="stats.totalDownloads"
            :value-style="{ color: '#1890ff' }"
          >
            <template #prefix>
              <DownloadOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :lg="8">
        <a-card title="活跃用户" :loading="loading">
          <a-statistic
            title="活跃用户"
            :value="stats.activeUsers"
            :value-style="{ color: '#52c41a' }"
          >
            <template #prefix>
              <CheckCircleOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="[16, 16]" style="margin-top: 16px">
      <a-col :xs="24" :lg="12">
        <a-card title="最近授权" :loading="loading">
          <a-list :data-source="recentLicenses" size="small">
            <template #renderItem="{ item }">
              <a-list-item>
                <a-list-item-meta>
                  <template #title>
                    <a-tag :color="getLicenseStatusColor(item.status)">
                      {{ getLicenseStatusText(item.status) }}
                    </a-tag>
                    {{ item.product?.name || '-' }}
                  </template>
                  <template #description>
                    {{ item.email }} - {{ formatTime(item.createdAt) }}
                  </template>
                </a-list-item-meta>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </a-col>
      <a-col :xs="24" :lg="12">
        <a-card title="最近插件" :loading="loading">
          <a-list :data-source="recentPlugins" size="small">
            <template #renderItem="{ item }">
              <a-list-item>
                <a-list-item-meta>
                  <template #title>
                    <a-tag :color="getPluginStatusColor(item.status)">
                      {{ getPluginStatusText(item.status) }}
                    </a-tag>
                    {{ item.name }}
                  </template>
                  <template #description>
                    {{ item.version }} - {{ formatTime(item.createdAt) }}
                  </template>
                </a-list-item-meta>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import {
  DatabaseOutlined,
  AppstoreOutlined,
  UserOutlined,
  DownloadOutlined,
  CheckCircleOutlined
} from '@ant-design/icons-vue'
import { dashboardAPI } from '../api'

const loading = ref(false)
const stats = reactive({
  totalProducts: 0,
  totalVersions: 0,
  totalPlugins: 0,
  totalLicenses: 0,
  pendingLicenses: 0,
  approvedLicenses: 0,
  rejectedLicenses: 0,
  totalUsers: 0,
  activeUsers: 0,
  totalDownloads: 0
})
const recentLicenses = ref([])
const recentPlugins = ref([])

const getLicenseStatusColor = (status) => {
  const colors = {
    pending: 'orange',
    approved: 'green',
    rejected: 'red'
  }
  return colors[status] || 'default'
}

const getLicenseStatusText = (status) => {
  const texts = {
    pending: '待审核',
    approved: '已批准',
    rejected: '已拒绝'
  }
  return texts[status] || status
}

const getPluginStatusColor = (status) => {
  const colors = {
    active: 'green',
    inactive: 'red'
  }
  return colors[status] || 'default'
}

const getPluginStatusText = (status) => {
  const texts = {
    active: '启用',
    inactive: '停用'
  }
  return texts[status] || status
}

const formatTime = (time) => {
  if (!time) return '-'
  const date = new Date(time)
  return date.toLocaleString('zh-CN')
}

const loadDashboardStats = async () => {
  loading.value = true
  try {
    const response = await dashboardAPI.getStats()
    if (response.data && response.data.result) {
      Object.assign(stats, response.data.result)
    }
  } catch (error) {
    console.error('Load dashboard stats error:', error)
    message.error('加载统计数据失败')
  } finally {
    loading.value = false
  }
}

const loadRecentLicenses = async () => {
  try {
    const response = await dashboardAPI.getRecentLicenses()
    if (response.data && response.data.result) {
      recentLicenses.value = response.data.result
    }
  } catch (error) {
    console.error('Load recent licenses error:', error)
  }
}

const loadRecentPlugins = async () => {
  try {
    const response = await dashboardAPI.getRecentPlugins()
    if (response.data && response.data.result) {
      recentPlugins.value = response.data.result
    }
  } catch (error) {
    console.error('Load recent plugins error:', error)
  }
}

onMounted(() => {
  loadDashboardStats()
  loadRecentLicenses()
  loadRecentPlugins()
})
</script>

<style scoped>
.dashboard {
  padding: 24px;
}

:deep(.ant-card) {
  border-radius: 8px;
}

:deep(.ant-statistic-title) {
  font-size: 14px;
  color: #666;
}

:deep(.ant-statistic-content) {
  font-size: 24px;
  font-weight: bold;
}
</style>
