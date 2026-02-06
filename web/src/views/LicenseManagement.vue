<template>
  <div>
    <a-page-header title="授权管理" />
    <a-table
      :columns="columns"
      :data-source="licenses"
      :loading="loading"
      :pagination="{ pageSize: 10 }"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusText(record.status) }}
          </a-tag>
        </template>
        <template v-if="column.key === 'action'">
          <a-space v-if="record.status === 'pending'">
            <a-button
              type="primary"
              size="small"
              @click="approveLicense(record.id)"
            >
              批准
            </a-button>
            <a-button
              danger
              size="small"
              @click="rejectLicense(record.id)"
            >
              拒绝
            </a-button>
          </a-space>
          <span v-else>-</span>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { licenseAPI } from '../api'

const licenses = ref([])
const loading = ref(false)

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id' },
  { title: '用户ID', dataIndex: 'user_id', key: 'user_id' },
  { title: '插件ID', dataIndex: 'plugin_id', key: 'plugin_id' },
  { title: '状态', key: 'status' },
  { title: '申请理由', dataIndex: 'reason', key: 'reason' },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
  { title: '操作', key: 'action' }
]

const getStatusColor = (status) => {
  const colors = {
    pending: 'orange',
    approved: 'green',
    rejected: 'red'
  }
  return colors[status] || 'default'
}

const getStatusText = (status) => {
  const texts = {
    pending: '待审核',
    approved: '已批准',
    rejected: '已拒绝'
  }
  return texts[status] || status
}

const loadLicenses = async () => {
  loading.value = true
  try {
    const response = await licenseAPI.getAll()
    licenses.value = response.data
  } catch (error) {
    message.error('加载授权列表失败')
  } finally {
    loading.value = false
  }
}

const approveLicense = async (id) => {
  try {
    await licenseAPI.approve(id)
    message.success('批准成功')
    loadLicenses()
  } catch (error) {
    message.error('批准失败')
  }
}

const rejectLicense = async (id) => {
  try {
    await licenseAPI.reject(id)
    message.success('拒绝成功')
    loadLicenses()
  } catch (error) {
    message.error('拒绝失败')
  }
}

onMounted(() => {
  loadLicenses()
})
</script>
