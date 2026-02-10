<template>
  <div class="license-type-management">
    <div class="header">
      <a-button type="primary" @click="showCreateModal">新增授权类型</a-button>
    </div>

    <a-table
      :columns="columns"
      :data-source="licenseTypes"
      :loading="loading"
      :pagination="pagination"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a-space>
            <a-button type="link" size="small" @click="handleEdit(record)">编辑</a-button>
            <a-popconfirm
              title="确定要删除这个授权类型吗？"
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
      :title="isEdit ? '编辑授权类型' : '新增授权类型'"
      @ok="handleSubmit"
      @cancel="handleCancel"
    >
      <a-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        layout="vertical"
      >
        <a-form-item label="授权类型名称" name="name">
          <a-input v-model:value="formData.name" placeholder="请输入授权类型名称" />
        </a-form-item>
        <a-form-item label="授权类型代码" name="code">
          <a-input v-model:value="formData.code" placeholder="请输入授权类型代码" />
        </a-form-item>
        <a-form-item label="是否为付费版本" name="isPaid">
          <a-switch v-model:checked="formData.isPaid" />
        </a-form-item>
        <a-form-item label="创建人" name="createdBy">
          <a-input v-model:value="formData.createdBy" placeholder="请输入创建人" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { licenseTypeAPI } from '../api'

const loading = ref(false)
const licenseTypes = ref([])
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
  { title: '授权类型名称', dataIndex: 'name', key: 'name' },
  { title: '授权类型代码', dataIndex: 'code', key: 'code' },
  { title: '是否付费', dataIndex: 'isPaid', key: 'isPaid', width: 100, customRender: (text) => text ? '是' : '否' },
  { title: '创建人', dataIndex: 'createdBy', key: 'createdBy' },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 180 },
  { title: '操作', key: 'action', width: 150, fixed: 'right' }
]

const formData = reactive({
  name: '',
  code: '',
  isPaid: false,
  createdBy: ''
})

const rules = {
  name: [{ required: true, message: '请输入授权类型名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入授权类型代码', trigger: 'blur' }]
}

const loadLicenseTypes = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      pageSize: pagination.pageSize
    }
    const response = await licenseTypeAPI.getAll(params)
    console.log('License types response:', response)
    if (response.data && response.data.result) {
      licenseTypes.value = response.data.result
      pagination.total = response.data.total || response.data.result.length
    } else {
      licenseTypes.value = []
      pagination.total = 0
    }
  } catch (error) {
    console.error('Load license types error:', error)
    message.error('加载授权类型列表失败')
  } finally {
    loading.value = false
  }
}

const showCreateModal = () => {
  isEdit.value = false
  currentEditId.value = null
  Object.assign(formData, {
    name: '',
    code: '',
    isPaid: false,
    createdBy: ''
  })
  modalVisible.value = true
}

const handleEdit = (record) => {
  isEdit.value = true
  currentEditId.value = record.id
  Object.assign(formData, {
    name: record.name,
    code: record.code,
    isPaid: record.isPaid,
    createdBy: record.createdBy
  })
  modalVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    
    if (isEdit.value) {
      await licenseTypeAPI.update(currentEditId.value, formData)
      message.success('更新授权类型成功')
    } else {
      await licenseTypeAPI.create(formData)
      message.success('创建授权类型成功')
    }
    
    modalVisible.value = false
    await loadLicenseTypes()
  } catch (error) {
    if (error.errorFields) {
      return
    }
    console.error('Submit error:', error)
    message.error(isEdit.value ? '更新授权类型失败' : '创建授权类型失败')
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
    await licenseTypeAPI.delete(id)
    message.success('删除授权类型成功')
    await loadLicenseTypes()
  } catch (error) {
    console.error('Delete error:', error)
    message.error('删除授权类型失败')
  }
}

onMounted(() => {
  loadLicenseTypes()
})
</script>

<style scoped>
.license-type-management {
  padding: 24px;
}

.header {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-bottom: 24px;
}
</style>