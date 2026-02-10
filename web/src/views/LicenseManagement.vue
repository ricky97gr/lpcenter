<template>
  <div class="license-management">
    <a-table
      :columns="columns"
      :data-source="licenses"
      :loading="loading"
      :pagination="pagination"
      row-key="id"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusText(record.status) }}
          </a-tag>
        </template>
        <template v-if="column.key === 'product'">
          {{ record.product?.name || '-' }}
        </template>
        <template v-if="column.key === 'licenseType'">
          <a-tag :color="getLicenseTypeColor(record.licenseType)">
            {{ getLicenseTypeText(record.licenseType) }}
          </a-tag>
        </template>
        <template v-if="column.key === 'licensePoints'">
          <a-tag color="blue">
            {{ record.licensePoints }}
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
          <a-space v-else>
            <a-button
              type="primary"
              size="small"
              @click="viewLicense(record.id, record.serialNumber)"
            >
              <template #icon>
                <EyeOutlined />
              </template>
              查看
            </a-button>
            <a-button
              type="primary"
              size="small"
              @click="downloadLicense(record.id, record.serialNumber)"
            >
              <template #icon>
                <DownloadOutlined />
              </template>
              下载
            </a-button>
            <a-button
              type="link"
              size="small"
              danger
              @click="deleteLicense(record.id)"
            >
              删除
            </a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="viewModalVisible"
      title="查看授权字符串"
      :footer="null"
      width="600px"
    >
      <a-form layout="vertical">
        <a-form-item label="序列号">
          <a-input :value="currentLicense.serialNumber" readonly />
        </a-form-item>
        <a-form-item label="授权类型">
          <a-input :value="getLicenseTypeText(currentLicense.licenseType)" readonly />
        </a-form-item>
        <a-form-item label="授权点数">
          <a-input :value="currentLicense.licensePoints" readonly />
        </a-form-item>
        <a-form-item label="授权字符串">
          <a-textarea
            :value="currentLicense.licenseString"
            :rows="8"
            readonly
          />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="copyLicenseString">
              <template #icon>
                <CopyOutlined />
              </template>
              复制授权字符串
            </a-button>
            <a-button @click="viewModalVisible = false">关闭</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { DownloadOutlined, EyeOutlined, CopyOutlined } from '@ant-design/icons-vue'
import { licenseAPI, licenseTypeAPI } from '../api'

const loading = ref(false)
const licenses = ref([])
const licenseTypes = ref([])
const viewModalVisible = ref(false)
const currentLicense = reactive({
  serialNumber: '',
  licenseType: '',
  licensePoints: 0,
  licenseString: ''
})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total) => `共 ${total} 条`
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '序列号', dataIndex: 'serialNumber', key: 'serialNumber' },
  { title: '产品', dataIndex: 'product', key: 'product' },
  { title: '授权类型', dataIndex: 'licenseType', key: 'licenseType', width: 120 },
  { title: '授权点数', dataIndex: 'licensePoints', key: 'licensePoints', width: 100 },
  { title: '有效期', dataIndex: 'expiryDate', key: 'expiryDate', width: 120 },
  { title: '状态', dataIndex: 'status', key: 'status', width: 100 },
  { title: '备注', dataIndex: 'remarks', key: 'remarks', ellipsis: true },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 180 },
  { title: '操作', key: 'action', width: 200, fixed: 'right' }
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

const getLicenseTypeColor = (licenseType) => {
  const colors = {
    normal: 'blue',
    pro: 'green',
    plus: 'orange',
    max: 'red',
    enterprise: 'purple'
  }
  return colors[licenseType] || 'default'
}

const getLicenseTypeText = (licenseType) => {
  // 首先从加载的授权类型中查找
  const found = licenseTypes.value.find(lt => lt.code === licenseType)
  if (found) {
    return found.name
  }
  
  // 作为备用，使用硬编码映射
  const texts = {
    normal: '普通版',
    pro: '专业版',
    plus: '增强版',
    max: '高级版',
    enterprise: '企业版'
  }
  return texts[licenseType] || licenseType
}

const loadLicenses = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      pageSize: pagination.pageSize
    }
    const response = await licenseAPI.getAll(params)
    console.log('Licenses response:', response)
    if (response.data && response.data.result) {
      licenses.value = response.data.result
      pagination.total = response.data.total || response.data.result.length
    } else {
      licenses.value = []
      pagination.total = 0
    }
  } catch (error) {
    console.error('Load licenses error:', error)
    message.error('加载授权列表失败')
  } finally {
    loading.value = false
  }
}

const loadLicenseTypes = async () => {
  try {
    const response = await licenseTypeAPI.getAll()
    console.log('License types response:', response)
    if (response.data && response.data.result) {
      licenseTypes.value = response.data.result
    } else {
      licenseTypes.value = []
    }
  } catch (error) {
    console.error('Load license types error:', error)
  }
}

const approveLicense = async (id) => {
  try {
    await licenseAPI.approve(id)
    message.success('批准成功')
    await loadLicenses()
  } catch (error) {
    console.error('Approve error:', error)
    message.error('批准失败')
  }
}

const rejectLicense = async (id) => {
  try {
    await licenseAPI.reject(id)
    message.success('拒绝成功')
    await loadLicenses()
  } catch (error) {
    console.error('Reject error:', error)
    message.error('拒绝失败')
  }
}

const deleteLicense = async (id) => {
  try {
    await licenseAPI.delete(id)
    message.success('删除成功')
    await loadLicenses()
  } catch (error) {
    console.error('Delete error:', error)
    message.error('删除失败')
  }
}

const downloadLicense = async (id, serialNumber) => {
  try {
    const response = await licenseAPI.download(id)
    
    const blob = new Blob([response.data], { type: 'application/octet-stream' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${serialNumber}.lic`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    message.success('下载成功')
  } catch (error) {
    console.error('Download error:', error)
    message.error('下载失败')
  }
}

const viewLicense = async (id, serialNumber) => {
  try {
    const response = await licenseAPI.getById(id)
    if (response.data && response.data.result) {
      Object.assign(currentLicense, {
        serialNumber: serialNumber,
        licenseType: response.data.result.licenseType || '',
        licensePoints: response.data.result.licensePoints || 0,
        licenseString: response.data.result.licenseString || ''
      })
      viewModalVisible.value = true
    }
  } catch (error) {
    console.error('View license error:', error)
    message.error('获取授权信息失败')
  }
}

const copyLicenseString = async () => {
  try {
    await navigator.clipboard.writeText(currentLicense.licenseString)
    message.success('复制成功')
  } catch (error) {
    console.error('Copy error:', error)
    message.error('复制失败')
  }
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadLicenses()
}

onMounted(() => {
  loadLicenses()
  loadLicenseTypes()
})
</script>

<style scoped>
.license-management {
  padding: 24px;
}
</style>