<template>
  <div class="user-management">
    <div class="header">
      <a-button type="primary" @click="showCreateModal">新增用户</a-button>
    </div>

    <a-table
      :columns="columns"
      :data-source="users"
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
        <template v-if="column.key === 'action'">
          <a-space>
            <a-button type="link" size="small" @click="handleEdit(record)">编辑</a-button>
            <a-button type="link" size="small" @click="showResetPasswordModal(record)">重置密码</a-button>
            <a-popconfirm
              title="确定要删除这个用户吗？"
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
      :title="isEdit ? '编辑用户' : '新增用户'"
      @ok="handleSubmit"
      @cancel="handleCancel"
    >
      <a-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        layout="vertical"
      >
        <a-form-item label="用户名" name="username">
          <a-input v-model:value="formData.username" placeholder="请输入用户名" />
        </a-form-item>
        <a-form-item label="邮箱" name="email">
          <a-input v-model:value="formData.email" placeholder="请输入邮箱" />
        </a-form-item>
        <a-form-item label="状态" name="status">
          <a-select v-model:value="formData.status" placeholder="请选择状态">
            <a-select-option value="active">启用</a-select-option>
            <a-select-option value="inactive">停用</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal
      v-model:open="resetPasswordVisible"
      title="重置密码"
      @ok="handleResetPassword"
      @cancel="resetPasswordVisible = false"
    >
      <a-form
        ref="resetPasswordFormRef"
        :model="resetPasswordData"
        :rules="resetPasswordRules"
        layout="vertical"
      >
        <a-form-item label="新密码" name="newPassword">
          <a-input-password v-model:value="resetPasswordData.newPassword" placeholder="请输入新密码" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { userAPI } from '../api'

const loading = ref(false)
const users = ref([])
const modalVisible = ref(false)
const isEdit = ref(false)
const currentEditId = ref(null)
const resetPasswordVisible = ref(false)
const resetPasswordUserId = ref(null)
const formRef = ref()
const resetPasswordFormRef = ref()

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total) => `共 ${total} 条`
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '用户名', dataIndex: 'username', key: 'username' },
  { title: '邮箱', dataIndex: 'email', key: 'email' },
  { title: '状态', dataIndex: 'status', key: 'status', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 180 },
  { title: '操作', key: 'action', width: 200, fixed: 'right' }
]

const formData = reactive({
  username: '',
  email: '',
  status: 'active'
})

const resetPasswordData = reactive({
  newPassword: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

const resetPasswordRules = {
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ]
}

const getStatusColor = (status) => {
  const colors = {
    active: 'green',
    inactive: 'red',
    disabled: 'gray'
  }
  return colors[status] || 'default'
}

const getStatusText = (status) => {
  const texts = {
    active: '启用',
    inactive: '停用',
    disabled: '禁用'
  }
  return texts[status] || status
}

const loadUsers = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      pageSize: pagination.pageSize
    }
    const response = await userAPI.getAll(params)
    console.log('Users response:', response)
    if (response.data && response.data.result) {
      users.value = response.data.result
      pagination.total = response.data.total || response.data.result.length
    } else {
      users.value = []
      pagination.total = 0
    }
  } catch (error) {
    console.error('Load users error:', error)
    message.error('加载用户列表失败')
  } finally {
    loading.value = false
  }
}

const showCreateModal = () => {
  isEdit.value = false
  currentEditId.value = null
  Object.assign(formData, {
    username: '',
    email: '',
    status: 'active'
  })
  modalVisible.value = true
}

const handleEdit = (record) => {
  isEdit.value = true
  currentEditId.value = record.id
  Object.assign(formData, {
    username: record.username,
    email: record.email,
    status: record.status
  })
  modalVisible.value = true
}

const showResetPasswordModal = (record) => {
  resetPasswordUserId.value = record.id
  resetPasswordData.newPassword = ''
  resetPasswordVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    
    if (isEdit.value) {
      await userAPI.update(currentEditId.value, formData)
      message.success('更新用户成功')
    } else {
      await userAPI.create(formData)
      message.success('创建用户成功')
    }
    
    modalVisible.value = false
    await loadUsers()
  } catch (error) {
    if (error.errorFields) {
      return
    }
    console.error('Submit error:', error)
    message.error(isEdit.value ? '更新用户失败' : '创建用户失败')
  } finally {
    loading.value = false
  }
}

const handleResetPassword = async () => {
  try {
    await resetPasswordFormRef.value.validate()
    await userAPI.resetPassword(resetPasswordUserId.value, { newPassword: resetPasswordData.newPassword })
    message.success('重置密码成功')
    resetPasswordVisible.value = false
  } catch (error) {
    if (error.errorFields) {
      return
    }
    console.error('Reset password error:', error)
    message.error('重置密码失败')
  }
}

const handleCancel = () => {
  modalVisible.value = false
  formRef.value?.resetFields()
}

const handleDelete = async (id) => {
  try {
    await userAPI.delete(id)
    message.success('删除用户成功')
    await loadUsers()
  } catch (error) {
    console.error('Delete error:', error)
    message.error('删除用户失败')
  }
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  loadUsers()
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.user-management {
  padding: 24px;
}

.header {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-bottom: 24px;
}
</style>
