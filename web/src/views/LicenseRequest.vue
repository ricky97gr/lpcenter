<template>
  <div class="license-request">
    <a-card title="申请授权">
      <a-form
        ref="formRef"
        :model="form"
        :rules="rules"
        layout="vertical"
        @finish="handleSubmit"
      >
        <a-form-item
          label="序列号"
          name="serialNumber"
        >
          <a-input
            v-model:value="form.serialNumber"
            placeholder="请输入序列号"
            style="width: 400px;"
          />
        </a-form-item>

        <a-form-item
          label="产品"
          name="productId"
        >
          <a-select
            v-model:value="form.productId"
            placeholder="请选择产品"
            style="width: 400px;"
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

        <a-form-item
          label="授权类型"
          name="version"
        >
          <a-select
            v-model:value="form.version"
            placeholder="请选择授权类型"
            style="width: 400px;"
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

        <a-form-item
          label="有效期"
          name="expiryDate"
        >
          <a-date-picker
            v-model:value="form.expiryDate"
            placeholder="请选择有效期"
            style="width: 400px;"
            :disabled-date="disabledDate"
          />
        </a-form-item>

        <a-form-item
          label="备注"
          name="remarks"
        >
          <a-textarea
            v-model:value="form.remarks"
            placeholder="请输入备注信息（可选）"
            :rows="4"
            style="width: 400px;"
          />
        </a-form-item>

        <a-form-item>
          <a-space>
            <a-button type="primary" html-type="submit" :loading="loading">
              <template #icon>
                <CheckOutlined />
              </template>
              提交申请
            </a-button>
            <a-button @click="handleReset">
              <template #icon>
                <ReloadOutlined />
              </template>
              重置
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { CheckOutlined, ReloadOutlined } from '@ant-design/icons-vue'
import { productAPI, licenseAPI, versionAPI } from '../api'
import dayjs from 'dayjs'

const loading = ref(false)
const productLoading = ref(false)
const versionLoading = ref(false)
const formRef = ref()
const products = ref([])
const versions = ref([])

const form = reactive({
  serialNumber: '',
  productId: undefined,
  version: '',
  expiryDate: null,
  remarks: ''
})

const rules = {
  serialNumber: [{ required: true, message: '请输入序列号', trigger: 'blur' }],
  productId: [{ required: true, message: '请选择产品', trigger: 'change' }],
  version: [{ required: true, message: '请输入版本号', trigger: 'blur' }],
  expiryDate: [{ required: true, message: '请选择有效期', trigger: 'change' }]
}

const disabledDate = (current) => {
  return current && current < dayjs().startOf('day')
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
      versions.value = response.data.result.filter(v => v.isPaid === true)
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

const handleSubmit = async () => {
  loading.value = true
  try {
    await licenseAPI.create({
      serialNumber: form.serialNumber,
      productId: form.productId,
      version: form.version,
      expiryDate: form.expiryDate ? form.expiryDate.format('YYYY-MM-DD') : '',
      remarks: form.remarks
    })

    message.success('申请提交成功')
    handleReset()
  } catch (error) {
    console.error('Submit error:', error)
    message.error('申请提交失败')
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  formRef.value?.resetFields()
  Object.assign(form, {
    serialNumber: '',
    productId: undefined,
    version: '',
    expiryDate: null,
    remarks: ''
  })
}

onMounted(() => {
  loadProducts()
  loadVersions()
})
</script>

<style scoped>
.license-request {
  padding: 24px;
}
</style>
