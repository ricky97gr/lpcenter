<template>
  <div class="plugin-management">
    <div class="header">
      <a-button type="primary" @click="showCreateModal">发布插件</a-button>
    </div>

    <a-table
      :columns="columns"
      :data-source="plugins"
      :loading="loading"
      :pagination="pagination"
      row-key="id"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'versionType'">
          <a-tag :color="getVersionTypeColor(record.versionType)">
            {{ getVersionTypeLabel(record.versionType) }}
          </a-tag>
        </template>
        <template v-if="column.key === 'status'">
          <a-tag :color="getStatusColor(record.status)">
            {{ getStatusText(record.status) }}
          </a-tag>
          <a-tag v-if="record.signed" color="blue" style="margin-left: 4px">已签名</a-tag>
        </template>
        <template v-if="column.key === 'product'">
          {{ record.product?.name || '-' }}
        </template>
        <template v-if="column.key === 'action'">
          <a-space>
            <a-button type="link" size="small" @click="handleViewDetail(record)">
              <template #icon>
                <EyeOutlined />
              </template>
              查看详情
            </a-button>
            <a-dropdown>
              <a-button type="link" size="small">
                操作 <DownOutlined />
              </a-button>
              <template #overlay>
                <a-menu>
                  <a-menu-item 
                    v-if="record.status === 'pending' && !record.signed" 
                    @click="handleSign(record.id)"
                  >
                    签名
                  </a-menu-item>
                  <a-menu-item 
                    v-if="record.status === 'signed'" 
                    @click="handlePublish(record.id)"
                  >
                    发布
                  </a-menu-item>
                  <a-menu-item 
                    v-if="record.status === 'pending' && !record.signed" 
                    @click="handlePublish(record.id)"
                  >
                    直接发布
                  </a-menu-item>
                  <a-menu-item 
                    v-if="record.status === 'published'" 
                    @click="handleDisable(record.id)"
                  >
                    停用
                  </a-menu-item>
                  <a-menu-item 
                    v-if="record.status === 'disabled'" 
                    @click="handlePublish(record.id)"
                  >
                    重新发布
                  </a-menu-item>
                  <a-menu-item @click="handleEdit(record)">编辑</a-menu-item>
                  <a-menu-item danger @click="handleDelete(record.id)">删除</a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalVisible"
      :title="isEdit ? '编辑插件' : '发布插件'"
      width="800px"
      @ok="handleSubmit"
      @cancel="handleCancel"
    >
      <a-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        layout="vertical"
      >
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="选择产品" name="productId">
              <a-select
                v-model:value="formData.productId"
                placeholder="请选择产品"
                :loading="productLoading"
                show-search
                :filter-option="filterOption"
              >
                <a-select-option
                  v-for="product in products"
                  :key="product.id"
                  :value="product.id"
                >
                  {{ product.name }} ({{ product.code }})
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="版本类型" name="versionType">
              <a-select 
                v-model:value="formData.versionType" 
                placeholder="请选择版本类型"
                :loading="licenseTypeLoading"
                show-search
                :filter-option="filterLicenseTypeOption"
              >
                <a-select-option
                  v-for="licenseType in licenseTypes"
                  :key="licenseType.code"
                  :value="licenseType.code"
                >
                  {{ licenseType.name }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="插件编号" name="code">
              <a-input v-model:value="formData.code" placeholder="请输入插件编号" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="插件名称" name="name">
              <a-input v-model:value="formData.name" placeholder="请输入插件名称" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="版本号" name="version">
              <a-input v-model:value="formData.version" placeholder="请输入版本号，如：1.0.0" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="作者" name="author">
              <a-input v-model:value="formData.author" placeholder="请输入作者名称" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="最低服务器版本" name="miniServerVersion">
              <a-input v-model:value="formData.miniServerVersion" placeholder="请输入最低服务器版本，如：1.0.0" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="最低客户端版本" name="miniClientVersion">
              <a-input v-model:value="formData.miniClientVersion" placeholder="请输入最低客户端版本，如：1.0.0" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="文件路径" name="filePath">
              <a-input v-model:value="formData.filePath" placeholder="请输入文件路径" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="下载地址" name="downloadUrl">
              <a-input v-model:value="formData.downloadUrl" placeholder="请输入下载地址（可选）" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="插件描述" name="description">
          <a-textarea v-model:value="formData.description" :rows="3" placeholder="请输入插件描述" />
        </a-form-item>
        <a-form-item label="提示信息" name="tips">
          <a-textarea v-model:value="formData.tips" :rows="3" placeholder="请输入提示信息" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal
      v-model:open="detailModalVisible"
      title="插件详情"
      :footer="null"
      width="800px"
    >
      <a-descriptions bordered :column="2">
        <a-descriptions-item label="插件名称">
          {{ currentPlugin.name }}
        </a-descriptions-item>
        <a-descriptions-item label="插件编号">
          {{ currentPlugin.code }}
        </a-descriptions-item>
        <a-descriptions-item label="版本号">
          {{ currentPlugin.version }}
        </a-descriptions-item>
        <a-descriptions-item label="授权类型">
          <a-tag :color="getVersionTypeColor(currentPlugin.versionType)">
            {{ getVersionTypeLabel(currentPlugin.versionType) }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="产品">
          {{ currentPlugin.product?.name || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="作者">
          {{ currentPlugin.author || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="最低服务器版本">
          {{ currentPlugin.miniServerVersion || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="最低客户端版本">
          {{ currentPlugin.miniClientVersion || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="下载量">
          {{ currentPlugin.downloadCount || 0 }}
        </a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag :color="getStatusColor(currentPlugin.status)">
            {{ getStatusText(currentPlugin.status) }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="文件路径" :span="2">
          {{ currentPlugin.filePath }}
        </a-descriptions-item>
        <a-descriptions-item label="下载地址" :span="2">
          {{ currentPlugin.downloadUrl || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="插件描述" :span="2">
          {{ currentPlugin.description || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="提示信息" :span="2">
          {{ currentPlugin.tips || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="上传时间" :span="2">
          {{ formatTime(currentPlugin.uploadedAt) }}
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { DownOutlined, EyeOutlined } from '@ant-design/icons-vue'
import { pluginAPI, productAPI, licenseTypeAPI } from '../api'

const loading = ref(false)
const productLoading = ref(false)
const licenseTypeLoading = ref(false)
const plugins = ref([])
const products = ref([])
const licenseTypes = ref([])
const modalVisible = ref(false)
const detailModalVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()
const currentEditId = ref(null)
const currentPlugin = ref({})

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total) => `共 ${total} 条`
})

const columns = [
  { title: '插件名称', dataIndex: 'name', key: 'name' },
  { title: '产品', dataIndex: 'product', key: 'product' },
  { title: '版本号', dataIndex: 'version', key: 'version', width: 100 },
  { title: '授权类型', dataIndex: 'versionType', key: 'versionType', width: 100 },
  { title: '作者', dataIndex: 'author', key: 'author', width: 100 },
  { title: '状态', dataIndex: 'status', key: 'status', width: 100 },
  { title: '下载量', dataIndex: 'downloadCount', key: 'downloadCount', width: 100 },
  { title: '操作', key: 'action', width: 300, fixed: 'right' }
]

const formData = reactive({
  productId: null,
  versionType: null,
  code: '',
  name: '',
  version: '',
  miniServerVersion: '',
  miniClientVersion: '',
  description: '',
  tips: '',
  author: '',
  filePath: '',
  downloadUrl: ''
})

const rules = {
  productId: [{ required: true, message: '请选择产品', trigger: 'change' }],
  versionType: [{ required: true, message: '请选择版本类型', trigger: 'change' }],
  code: [{ required: true, message: '请输入插件编号', trigger: 'blur' }],
  name: [{ required: true, message: '请输入插件名称', trigger: 'blur' }],
  version: [{ required: true, message: '请输入版本号', trigger: 'blur' }],
  filePath: [{ required: true, message: '请输入文件路径', trigger: 'blur' }]
}

const getVersionTypeColor = (type) => {
  const colors = {
    normal: 'blue',
    pro: 'green',
    plus: 'orange',
    max: 'red'
  }
  return colors[type] || 'default'
}

const getVersionTypeLabel = (type) => {
  const found = licenseTypes.value.find(v => v.code === type)
  return found ? found.name : type
}

const getStatusColor = (status) => {
  const colors = {
    pending: 'orange',
    published: 'green',
    signed: 'blue',
    disabled: 'red'
  }
  return colors[status] || 'default'
}

const getStatusText = (status) => {
  const texts = {
    pending: '待审核',
    published: '已发布',
    signed: '已签名',
    disabled: '已停用'
  }
  return texts[status] || status
}

const filterOption = (input, option) => {
  const product = products.value.find(p => p.id === option.value)
  return product && (product.name.toLowerCase().includes(input.toLowerCase()) || 
                     product.code.toLowerCase().includes(input.toLowerCase()))
}

const filterLicenseTypeOption = (input, option) => {
  const licenseType = licenseTypes.value.find(v => v.code === option.value)
  return licenseType && (licenseType.name.toLowerCase().includes(input.toLowerCase()) ||
                     licenseType.code.toLowerCase().includes(input.toLowerCase()))
}

const loadPlugins = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      pageSize: pagination.pageSize
    }
    const response = await pluginAPI.getAll(params)
    console.log('Plugins response:', response)
    if (response.data && response.data.result) {
      plugins.value = response.data.result
      pagination.total = response.data.total || response.data.result.length
    } else {
      plugins.value = []
      pagination.total = 0
    }
  } catch (error) {
    console.error('Load plugins error:', error)
    message.error('加载插件列表失败')
  } finally {
    loading.value = false
  }
}

const loadProducts = async () => {
  productLoading.value = true
  try {
    const response = await productAPI.getAll()
    console.log('Products response:', response)
    if (response.data && response.data.result) {
      products.value = response.data.result
    } else {
      products.value = []
    }
  } catch (error) {
    console.error('Load products error:', error)
    message.error('加载产品列表失败')
  } finally {
    productLoading.value = false
  }
}

const loadLicenseTypes = async () => {
  licenseTypeLoading.value = true
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
    message.error('加载授权类型列表失败')
  } finally {
    licenseTypeLoading.value = false
  }
}

const showCreateModal = () => {
  isEdit.value = false
  currentEditId.value = null
  Object.assign(formData, {
    productId: null,
    versionType: null,
    code: '',
    name: '',
    version: '',
    miniServerVersion: '',
    miniClientVersion: '',
    description: '',
    tips: '',
    author: '',
    filePath: '',
    downloadUrl: ''
  })
  modalVisible.value = true
}

const handleEdit = (record) => {
  isEdit.value = true
  currentEditId.value = record.id
  Object.assign(formData, {
    productId: record.productId,
    versionType: record.versionType,
    code: record.code,
    name: record.name,
    version: record.version,
    miniServerVersion: record.miniServerVersion,
    miniClientVersion: record.miniClientVersion,
    description: record.description,
    tips: record.tips,
    author: record.author,
    filePath: record.filePath,
    downloadURL: record.downloadUrl
  })
  modalVisible.value = true
}

const handleViewDetail = (record) => {
  currentPlugin.value = { ...record }
  detailModalVisible.value = true
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    
    if (isEdit.value) {
      await pluginAPI.update(currentEditId.value, formData)
      message.success('更新插件成功')
    } else {
      await pluginAPI.create(formData)
      message.success('发布插件成功')
    }
    
    modalVisible.value = false
    await loadPlugins()
  } catch (error) {
    if (error.errorFields) {
      return
    }
    console.error('Submit error:', error)
    message.error(isEdit.value ? '更新插件失败' : '发布插件失败')
  } finally {
    loading.value = false
  }
}

const handleCancel = () => {
  modalVisible.value = false
  formRef.value?.resetFields()
}

const handleDelete = async (id) => {
  try {
    await pluginAPI.delete(id)
    message.success('删除插件成功')
    await loadPlugins()
  } catch (error) {
    console.error('Delete error:', error)
    message.error('删除插件失败')
  }
}

const handleApprove = async (id) => {
  try {
    await pluginAPI.updateStatus(id, { status: 'published' })
    message.success('插件已发布')
    await loadPlugins()
  } catch (error) {
    console.error('Approve error:', error)
    message.error('发布失败')
  }
}

const handleReject = async (id) => {
  try {
    await pluginAPI.updateStatus(id, { status: 'disabled' })
    message.success('插件已停用')
    await loadPlugins()
  } catch (error) {
    console.error('Reject error:', error)
    message.error('停用失败')
  }
}

const handleSign = async (id) => {
  try {
    await pluginAPI.sign(id)
    message.success('插件签名成功')
    await loadPlugins()
  } catch (error) {
    console.error('Sign error:', error)
    message.error('签名失败')
  }
}

const handlePublish = async (id) => {
  try {
    await pluginAPI.publish(id)
    message.success('插件发布成功')
    await loadPlugins()
  } catch (error) {
    console.error('Publish error:', error)
    message.error('发布失败')
  }
}

const handleDisable = async (id) => {
  try {
    await pluginAPI.disable(id)
    message.success('插件已停用')
    await loadPlugins()
  } catch (error) {
    console.error('Disable error:', error)
    message.error('停用失败')
  }
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadPlugins()
}

onMounted(() => {
  loadPlugins()
  loadProducts()
  loadLicenseTypes()
})
</script>

<style scoped>
.plugin-management {
  padding: 24px;
}

.header {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-bottom: 24px;
}
</style>