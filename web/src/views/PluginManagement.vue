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
        </template>
        <template v-if="column.key === 'product'">
          {{ record.product?.name || '-' }}
        </template>
        <template v-if="column.key === 'action'">
          <a-space>
            <a-dropdown v-if="record.status === 'pending'">
              <a-button type="link" size="small">
                审核 <DownOutlined />
              </a-button>
              <template #overlay>
                <a-menu>
                  <a-menu-item @click="handleApprove(record.id)">发布</a-menu-item>
                  <a-menu-item @click="handleReject(record.id)">停用</a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
            <a-button v-if="record.status === 'published'" type="link" size="small" @click="handleReject(record.id)">停用</a-button>
            <a-button v-if="record.status === 'disabled'" type="link" size="small" @click="handleApprove(record.id)">重新发布</a-button>
            <a-button type="link" size="small" @click="handleEdit(record)">编辑</a-button>
            <a-popconfirm
              title="确定要删除这个插件吗？"
              @confirm="handleDelete(record.id)"
            >
              <a-button type="link" size="small" danger>删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalVisible"
      :title="isEdit ? '编辑插件' : '发布插件'"
      width="600px"
      @ok="handleSubmit"
      @cancel="handleCancel"
    >
      <a-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        layout="vertical"
      >
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
        <a-form-item label="版本类型" name="versionType">
          <a-select 
            v-model:value="formData.versionType" 
            placeholder="请选择版本类型"
            :loading="versionLoading"
            show-search
            :filter-option="filterVersionOption"
          >
            <a-select-option
              v-for="version in versions"
              :key="version.code"
              :value="version.code"
            >
              {{ version.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="插件编号" name="code">
          <a-input v-model:value="formData.code" placeholder="请输入插件编号" />
        </a-form-item>
        <a-form-item label="版本号" name="version">
          <a-input v-model:value="formData.version" placeholder="请输入版本号，如：1.0.0" />
        </a-form-item>
        <a-form-item label="插件描述" name="description">
          <a-textarea v-model:value="formData.description" :rows="3" placeholder="请输入插件描述" />
        </a-form-item>
        <a-form-item label="作者" name="author">
          <a-input v-model:value="formData.author" placeholder="请输入作者名称" />
        </a-form-item>
        <a-form-item label="文件路径" name="filePath">
          <a-input v-model:value="formData.filePath" placeholder="请输入文件路径" />
        </a-form-item>
        <a-form-item label="下载地址" name="downloadUrl">
          <a-input v-model:value="formData.downloadUrl" placeholder="请输入下载地址（可选）" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { DownOutlined } from '@ant-design/icons-vue'
import { pluginAPI, productAPI, versionAPI } from '../api'

const loading = ref(false)
const productLoading = ref(false)
const versionLoading = ref(false)
const plugins = ref([])
const products = ref([])
const versions = ref([])
const modalVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()
const currentEditId = ref(null)

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total) => `共 ${total} 条`
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '插件编号', dataIndex: 'code', key: 'code' },
  { title: '版本号', dataIndex: 'version', key: 'version', width: 100 },
  { title: '版本类型', dataIndex: 'versionType', key: 'versionType', width: 100 },
  { title: '产品', dataIndex: 'product', key: 'product' },
  { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
  { title: '作者', dataIndex: 'author', key: 'author', width: 100 },
  { title: '状态', dataIndex: 'status', key: 'status', width: 100 },
  { title: '下载量', dataIndex: 'downloadCount', key: 'downloadCount', width: 100 },
  { title: '上传时间', dataIndex: 'uploadedAt', key: 'uploadedAt', width: 180 },
  { title: '操作', key: 'action', width: 200, fixed: 'right' }
]

const formData = reactive({
  productId: null,
  versionType: null,
  code: '',
  version: '',
  description: '',
  author: '',
  filePath: '',
  downloadUrl: ''
})

const rules = {
  productId: [{ required: true, message: '请选择产品', trigger: 'change' }],
  versionType: [{ required: true, message: '请选择版本类型', trigger: 'change' }],
  code: [{ required: true, message: '请输入插件编号', trigger: 'blur' }],
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
  const found = versions.value.find(v => v.code === type)
  return found ? found.name : type
}

const getStatusColor = (status) => {
  const colors = {
    pending: 'orange',
    published: 'green',
    disabled: 'red'
  }
  return colors[status] || 'default'
}

const getStatusText = (status) => {
  const texts = {
    pending: '审核中',
    published: '发布',
    disabled: '停用'
  }
  return texts[status] || status
}

const filterOption = (input, option) => {
  const product = products.value.find(p => p.id === option.value)
  return product && (product.name.toLowerCase().includes(input.toLowerCase()) || 
                     product.code.toLowerCase().includes(input.toLowerCase()))
}

const filterVersionOption = (input, option) => {
  const version = versions.value.find(v => v.code === option.value)
  return version && (version.name.toLowerCase().includes(input.toLowerCase()) ||
                     version.code.toLowerCase().includes(input.toLowerCase()))
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

const loadVersions = async () => {
  versionLoading.value = true
  try {
    const response = await versionAPI.getAll()
    console.log('Versions response:', response)
    if (response.data && response.data.result) {
      versions.value = response.data.result
    } else {
      versions.value = []
    }
  } catch (error) {
    console.error('Load versions error:', error)
    message.error('加载版本列表失败')
  } finally {
    versionLoading.value = false
  }
}

const showCreateModal = () => {
  isEdit.value = false
  currentEditId.value = null
  Object.assign(formData, {
    productId: null,
    versionType: null,
    code: '',
    version: '',
    description: '',
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
    version: record.version,
    description: record.description,
    author: record.author,
    filePath: record.filePath,
    downloadUrl: record.downloadUrl
  })
  modalVisible.value = true
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

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadPlugins()
}

onMounted(() => {
  loadPlugins()
  loadProducts()
  loadVersions()
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
