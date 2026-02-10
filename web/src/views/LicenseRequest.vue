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
          name="licenseType"
        >
          <a-select
            v-model:value="form.licenseType"
            placeholder="请选择授权类型"
            style="width: 400px;"
            :loading="licenseTypeLoading"
            show-search
            :filter-option="filterLicenseTypeOption"
          >
            <a-select-option
              v-for="licenseType in licenseTypes"
              :key="licenseType.id"
              :value="licenseType.code"
            >
              {{ licenseType.name }}
            </a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item
          label="授权点数"
          name="licensePoints"
        >
          <a-input-number
            v-model:value="form.licensePoints"
            placeholder="请输入授权点数"
            style="width: 400px;"
            :min="1"
            :max="99999"
            :step="1"
          />
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
import { productAPI, licenseAPI, licenseTypeAPI } from '../api'
import dayjs from 'dayjs'

const loading = ref(false)
const productLoading = ref(false)
const licenseTypeLoading = ref(false)
const formRef = ref()
const products = ref([])
const licenseTypes = ref([])

const form = reactive({
  serialNumber: '',
  productId: undefined,
  licenseType: '',
  licensePoints: 1,
  expiryDate: null,
  remarks: ''
})

const rules = {
  serialNumber: [{ required: true, message: '请输入序列号', trigger: 'blur' }],
  productId: [{ required: true, message: '请选择产品', trigger: 'change' }],
  licenseType: [{ required: true, message: '请选择授权类型', trigger: 'change' }],
  licensePoints: [{ required: true, message: '请输入授权点数', trigger: 'blur' }],
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

const filterLicenseTypeOption = (input, option) => {
  const licenseType = licenseTypes.value.find(v => v.code === option.value)
  return licenseType && (licenseType.name.toLowerCase().includes(input.toLowerCase()) ||
                         licenseType.code.toLowerCase().includes(input.toLowerCase()))
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
      licenseTypes.value = response.data.result.filter(v => v.isPaid === true)
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

const handleSubmit = async () => {
  loading.value = true
  try {
    await licenseAPI.create({
      serialNumber: form.serialNumber,
      productId: form.productId,
      licenseType: form.licenseType,
      licensePoints: form.licensePoints,
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
    licenseType: '',
    licensePoints: 1,
    expiryDate: null,
    remarks: ''
  })
}

onMounted(() => {
  loadProducts()
  loadLicenseTypes()
})
</script>

<style scoped>
.license-request {
  padding: 24px;
}
</style>