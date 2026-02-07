<template>
  <div class="version-management">
    <div class="header">
      <a-button type="primary" @click="showCreateModal">新增授权类型</a-button>
    </div>

    <a-table
      :columns="columns"
      :data-source="versions"
      :loading="loading"
      :pagination="pagination"
      row-key="id"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'isPaid'">
          <a-switch
            :checked="record.isPaid"
            @change="(checked) => handleTogglePaid(record.id, checked)"
          />
        </template>
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
          <a-input v-model:value="formData.name" placeholder="请输入授权类型名称，如：普通版" />
        </a-form-item>
        <a-form-item label="是否付费" name="isPaid">
          <a-switch v-model:checked="formData.isPaid" />
        </a-form-item>
        <a-form-item label="创建人" name="createdBy">
          <a-input v-model:value="formData.createdBy" placeholder="请输入创建人姓名（可选）" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { versionAPI } from '../api'

const loading = ref(false)
const versions = ref([])
const modalVisible = ref(false)
const isEdit = ref(false)
const currentEditId = ref(null)
const formRef = ref()

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
  { title: '是否付费', dataIndex: 'isPaid', key: 'isPaid', width: 100 },
  { title: '创建人', dataIndex: 'createdBy', key: 'createdBy', width: 120 },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 180 },
  { title: '操作', key: 'action', width: 150, fixed: 'right' }
]

const formData = reactive({
  name: '',
  isPaid: false,
  createdBy: ''
})

const rules = {
  name: [{ required: true, message: '请输入版本名称', trigger: 'blur' }]
}

const loadVersions = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      pageSize: pagination.pageSize
    }
    const response = await versionAPI.getAll(params)
    console.log('Versions response:', response)
    if (response.data && response.data.result) {
      versions.value = response.data.result
      pagination.total = response.data.total || response.data.result.length
    } else {
      versions.value = []
      pagination.total = 0
    }
  } catch (error) {
    console.error('Load versions error:', error)
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
    isPaid: record.isPaid,
    createdBy: record.createdBy
  })
  modalVisible.value = true
}

const handleTogglePaid = async (id, checked) => {
  try {
    const version = versions.value.find(v => v.id === id)
    if (version) {
      await versionAPI.update(id, { name: version.name, isPaid: checked, createdBy: version.createdBy })
      version.isPaid = checked
      message.success(checked ? '已设置为付费版本' : '已设置为免费版本')
    }
  } catch (error) {
    console.error('Toggle paid error:', error)
    message.error('更新失败')
  }
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    
    if (isEdit.value) {
      await versionAPI.update(currentEditId.value, formData)
      message.success('更新版本成功')
    } else {
      await versionAPI.create(formData)
      message.success('创建版本成功')
    }
    
    modalVisible.value = false
    await loadVersions()
  } catch (error) {
    if (error.errorFields) {
      return
    }
    console.error('Submit error:', error)
    message.error(isEdit.value ? '更新版本失败' : '创建版本失败')
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
    await versionAPI.delete(id)
    message.success('删除版本成功')
    await loadVersions()
  } catch (error) {
    console.error('Delete error:', error)
    message.error('删除版本失败')
  }
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadVersions()
}

onMounted(() => {
  loadVersions()
})
</script>

<style scoped>
.version-management {
  padding: 24px;
}

.header {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-bottom: 24px;
}
</style>
