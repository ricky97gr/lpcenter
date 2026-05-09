<template>
  <div class="plugin-upload">
    <div class="upload-container">
      <div class="upload-header">
        <h1>插件上传</h1>
        <p class="subtitle">上传您的插件，让更多用户使用</p>
      </div>

      <a-card class="upload-card">
        <a-form
          ref="formRef"
          :model="formData"
          :rules="rules"
          layout="vertical"
          @finish="handleSubmit"
        >
          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item label="选择产品" name="productUuid">
                <a-select
                  v-model:value="formData.productUuid"
                  placeholder="请选择产品"
                  :loading="productLoading"
                  show-search
                  :filter-option="filterOption"
                  size="large"
                >
                  <a-select-option
                  v-for="product in products"
                  :key="product.uuid"
                  :value="product.uuid"
                >
                  {{ product.name }} ({{ product.code }})
                </a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="授权类型" name="licenseType">
                <a-select 
                  v-model:value="formData.licenseType" 
                  placeholder="请选择授权类型"
                  :loading="licenseLoading"
                  show-search
                  :filter-option="filterLicenseTypeOption"
                  size="large"
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

          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item label="插件名称" name="name">
                <a-input 
                  v-model:value="formData.name" 
                  placeholder="请输入插件名称"
                  size="large"
                />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="联系方式" name="contact">
                <a-input 
                  v-model:value="formData.contact" 
                  placeholder="请输入联系方式（邮箱或电话）"
                  size="large"
                />
              </a-form-item>
            </a-col>
          </a-row>

          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item label="版本号" name="version">
                <a-input 
                  v-model:value="formData.version" 
                  placeholder="请输入版本号，如：1.0.0"
                  size="large"
                />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="作者" name="author">
                <a-input 
                  v-model:value="formData.author" 
                  placeholder="请输入作者名称"
                  size="large"
                />
              </a-form-item>
            </a-col>
          </a-row>

          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item label="最低服务器版本" name="miniServerVersion">
                <a-input 
                  v-model:value="formData.miniServerVersion" 
                  placeholder="请输入最低服务器版本，如：1.0.0"
                  size="large"
                />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="最低客户端版本" name="miniClientVersion">
                <a-input 
                  v-model:value="formData.miniClientVersion" 
                  placeholder="请输入最低客户端版本，如：1.0.0"
                  size="large"
                />
              </a-form-item>
            </a-col>
          </a-row>



          <a-form-item label="插件描述" name="description">
            <a-textarea 
              v-model:value="formData.description" 
              :rows="4" 
              placeholder="请输入插件描述，介绍插件的功能和特点"
            />
          </a-form-item>

          <a-form-item label="提示信息" name="tips">
            <a-textarea 
              v-model:value="formData.tips" 
              :rows="3" 
              placeholder="请输入提示信息，如使用注意事项、安装说明等"
            />
          </a-form-item>

          <a-form-item label="插件文件" name="file">
            <a-upload
              v-model:file-list="fileList"
              :before-upload="beforeUpload"
              :remove="handleRemove"
              :max-count="1"
              accept=".zip"
            >
              <a-button size="large">
                <template #icon>
                  <UploadOutlined />
                </template>
                选择ZIP文件
              </a-button>
            </a-upload>
            <div class="upload-tip">
              <a-typography-text type="secondary">
                仅支持ZIP格式文件，文件大小不超过100MB
              </a-typography-text>
            </div>
          </a-form-item>

          <a-form-item>
            <a-space size="large">
              <a-button 
                type="primary" 
                html-type="submit" 
                size="large"
                :loading="submitting"
                block
              >
                提交插件
              </a-button>
              <a-button size="large" @click="handleReset" block>
                重置表单
              </a-button>
            </a-space>
          </a-form-item>
        </a-form>
      </a-card>

      <a-card class="info-card">
        <a-typography-title :level="5">上传须知</a-typography-title>
        <a-list :data-source="uploadTips" size="small">
          <template #renderItem="{ item }">
            <a-list-item>
              <template #prefix>
                <CheckCircleOutlined style="color: #52c41a" />
              </template>
              {{ item }}
            </a-list-item>
          </template>
        </a-list>
      </a-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { UploadOutlined, CheckCircleOutlined } from '@ant-design/icons-vue'
import { publicAPI } from '../api'

const formRef = ref()
const productLoading = ref(false)
const licenseLoading = ref(false)
const submitting = ref(false)
const products = ref([])
const licenseTypes = ref([])
const fileList = ref([])

const formData = reactive({
  productUuid: null,
  licenseType: null,
  name: '',
  version: '',
  miniServerVersion: '',
  miniClientVersion: '',
  description: '',
  tips: '',
  author: '',
  contact: ''
})

const rules = {
  productUuid: [{ required: true, message: '请选择产品', trigger: 'change' }],
  licenseType: [{ required: true, message: '请选择授权类型', trigger: 'change' }],
  name: [{ required: true, message: '请输入插件名称', trigger: 'blur' }],
  version: [
    { required: true, message: '请输入版本号', trigger: 'blur' },
    { pattern: /^\d+\.\d+\.\d+$/, message: '版本号格式不正确，应为 x.y.z 格式', trigger: 'blur' }
  ],
  author: [{ required: true, message: '请输入作者名称', trigger: 'blur' }],
  contact: [
    { required: true, message: '请输入联系方式', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  description: [{ required: true, message: '请输入插件描述', trigger: 'blur' }]
}

const uploadTips = [
  '插件文件必须是ZIP格式',
  '文件大小不超过100MB',
  '插件上传后需要经过审核才能发布',
  '审核通过后，插件将自动签名并发布',
  '请确保插件内容合法合规，不包含恶意代码'
]

const filterOption = (input, option) => {
  const product = products.value.find(p => p.uuid === option.value)
  return product && (product.name.toLowerCase().includes(input.toLowerCase()) || 
                     product.code.toLowerCase().includes(input.toLowerCase()))
}

const filterLicenseTypeOption = (input, option) => {
  const licenseType = licenseTypes.value.find(v => v.code === option.value)
  return licenseType && (licenseType.name.toLowerCase().includes(input.toLowerCase()) ||
                     licenseType.code.toLowerCase().includes(input.toLowerCase()))
}

const loadProducts = async () => {
  productLoading.value = true
  try {
    const response = await publicAPI.getProducts()
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
  licenseLoading.value = true
  try {
    const response = await publicAPI.getLicenseTypes()
    if (response.data && response.data.result) {
      licenseTypes.value = response.data.result
    } else {
      licenseTypes.value = []
    }
  } catch (error) {
    console.error('Load license types error:', error)
    message.error('加载授权类型列表失败')
  } finally {
    licenseLoading.value = false
  }
}

const beforeUpload = (file) => {
  const isZip = file.type === 'application/zip' || file.name.endsWith('.zip')
  if (!isZip) {
    message.error('只能上传ZIP格式的文件')
    return false
  }
  const isLt100M = file.size / 1024 / 1024 < 100
  if (!isLt100M) {
    message.error('文件大小不能超过100MB')
    return false
  }
  fileList.value = [file]
  return false
}

const handleRemove = () => {
  fileList.value = []
  formData.filePath = ''
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    
    if (fileList.value.length === 0) {
      message.error('请选择要上传的插件文件')
      return
    }

    submitting.value = true
    
    const uploadFormData = new FormData()
    uploadFormData.append('productUuid', formData.productUuid)
    uploadFormData.append('licenseType', formData.licenseType)
    uploadFormData.append('name', formData.name)
    uploadFormData.append('version', formData.version)
    uploadFormData.append('miniServerVersion', formData.miniServerVersion)
    uploadFormData.append('miniClientVersion', formData.miniClientVersion)
    uploadFormData.append('description', formData.description)
    uploadFormData.append('tips', formData.tips)
    uploadFormData.append('author', formData.author)
    uploadFormData.append('contact', formData.contact)
    uploadFormData.append('file', fileList.value[0].originFileObj)

    await publicAPI.uploadPlugin(uploadFormData)
    message.success('插件上传成功，等待审核')
    handleReset()
  } catch (error) {
    if (error.errorFields) {
      return
    }
    console.error('Submit error:', error)
    message.error('插件上传失败')
  } finally {
    submitting.value = false
  }
}

const handleReset = () => {
  formRef.value?.resetFields()
  fileList.value = []
  Object.assign(formData, {
    productUuid: null,
    licenseType: null,
    code: '',
    name: '',
    version: '',
    miniServerVersion: '',
    miniClientVersion: '',
    description: '',
    tips: '',
    author: '',
    contact: '',
    filePath: ''
  })
}

onMounted(() => {
  loadProducts()
  loadLicenseTypes()
})
</script>

<style scoped>
.plugin-upload {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 40px 20px;
}

.upload-container {
  max-width: 1000px;
  margin: 0 auto;
}

.upload-header {
  text-align: center;
  color: white;
  margin-bottom: 40px;
}

.upload-header h1 {
  font-size: 48px;
  font-weight: 700;
  margin-bottom: 12px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.subtitle {
  font-size: 18px;
  opacity: 0.9;
  margin: 0;
}

.upload-card {
  margin-bottom: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  overflow: hidden;
}

.upload-card :deep(.ant-card-body) {
  padding: 15px !important;
}

.info-card {
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  overflow: hidden;
}

.info-card :deep(.ant-card-body) {
  padding: 15px !important;
}

.upload-tip {
  margin-top: 8px;
}

:deep(.ant-form-item-label > label) {
  font-weight: 500;
  font-size: 14px;
}

:deep(.ant-upload-list) {
  margin-top: 8px;
}

:deep(.ant-list-item) {
  padding: 8px 0;
}

@media (max-width: 768px) {
  .plugin-upload {
    padding: 20px 10px;
  }

  .upload-header h1 {
    font-size: 32px;
  }

  .subtitle {
    font-size: 16px;
  }

  .upload-card :deep(.ant-card-body) {
    padding: 20px;
  }
}
</style>
