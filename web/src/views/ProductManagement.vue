<template>
  <div class="product-management">
    <div class="header">
      <a-button type="primary" @click="showCreateModal">新增产品</a-button>
    </div>

    <a-table
      :columns="columns"
      :data-source="products"
      :loading="loading"
      :pagination="pagination"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="record.status === 'active' ? 'green' : 'red'">
            {{ record.status === 'active' ? '启用' : '禁用' }}
          </a-tag>
        </template>
        <template v-if="column.key === 'action'">
          <a-space>
            <a-button type="link" size="small" @click="handleEdit(record)">编辑</a-button>
            <a-popconfirm
              title="确定要删除这个产品吗？"
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
      :title="isEdit ? '编辑产品' : '新增产品'"
      @ok="handleSubmit"
      @cancel="handleCancel"
    >
      <a-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        layout="vertical"
      >
        <a-form-item label="产品名称" name="name">
          <a-input v-model:value="formData.name" placeholder="请输入产品名称" />
        </a-form-item>
        <a-form-item label="产品编号" name="code">
          <a-input v-model:value="formData.code" placeholder="请输入产品编号" />
        </a-form-item>
        <a-form-item label="产品描述" name="description">
          <a-textarea v-model:value="formData.description" :rows="4" placeholder="请输入产品描述" />
        </a-form-item>
        <a-form-item label="状态" name="status">
          <a-select v-model:value="formData.status" placeholder="请选择状态">
            <a-select-option value="active">启用</a-select-option>
            <a-select-option value="inactive">禁用</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { productAPI } from '../api'

const loading = ref(false)
const products = ref([])
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
  { title: '产品名称', dataIndex: 'name', key: 'name' },
  { title: '产品编号', dataIndex: 'code', key: 'code' },
  { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
  { title: '状态', dataIndex: 'status', key: 'status', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 180 },
  { title: '操作', key: 'action', width: 150, fixed: 'right' }
]

const formData = reactive({
  name: '',
  code: '',
  description: '',
  status: 'active'
})

const rules = {
  name: [{ required: true, message: '请输入产品名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入产品编号', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

const loadProducts = async () => {
  loading.value = true
  try {
    const response = await productAPI.getAll()
    console.log('Products response:', response)
    if (response.data && response.data.result) {
      products.value = response.data.result
      pagination.total = response.data.total || response.data.result.length
    } else {
      products.value = []
      pagination.total = 0
    }
  } catch (error) {
    console.error('Load products error:', error)
    message.error('加载产品列表失败')
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
    description: '',
    status: 'active'
  })
  modalVisible.value = true
}

const handleEdit = (record) => {
  isEdit.value = true
  currentEditId.value = record.id
  Object.assign(formData, {
    name: record.name,
    code: record.code,
    description: record.description,
    status: record.status
  })
  modalVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    
    if (isEdit.value) {
      await productAPI.update(currentEditId.value, formData)
      message.success('更新产品成功')
    } else {
      await productAPI.create(formData)
      message.success('创建产品成功')
    }
    
    modalVisible.value = false
    await loadProducts()
  } catch (error) {
    if (error.errorFields) {
      return
    }
    console.error('Submit error:', error)
    message.error(isEdit.value ? '更新产品失败' : '创建产品失败')
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
    await productAPI.delete(id)
    message.success('删除产品成功')
    await loadProducts()
  } catch (error) {
    console.error('Delete error:', error)
    message.error('删除产品失败')
  }
}

onMounted(() => {
  loadProducts()
})
</script>

<style scoped>
.product-management {
  padding: 24px;
}

.header {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-bottom: 24px;
}
</style>
