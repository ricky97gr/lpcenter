<template>
  <div>
    <a-page-header title="插件列表" />
    <a-button type="primary" @click="showModal" style="margin-bottom: 16px">
      发布新插件
    </a-button>
    <a-table :columns="columns" :data-source="plugins" :loading="loading">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a-space>
            <a-button type="link" @click="viewDetail(record.id)">查看</a-button>
            <a-button type="link" @click="downloadPlugin(record)">下载</a-button>
            <a-button type="link" danger @click="deletePlugin(record.id)">删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="visible"
      title="发布新插件"
      @ok="handleOk"
      @cancel="handleCancel"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="插件名称">
          <a-input v-model:value="form.name" />
        </a-form-item>
        <a-form-item label="版本">
          <a-input v-model:value="form.version" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model:value="form.description" />
        </a-form-item>
        <a-form-item label="作者">
          <a-input v-model:value="form.author" />
        </a-form-item>
        <a-form-item label="文件路径">
          <a-input v-model:value="form.file_path" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { pluginAPI } from '../api'

const router = useRouter()
const plugins = ref([])
const loading = ref(false)
const visible = ref(false)
const form = ref({
  name: '',
  version: '',
  description: '',
  author: '',
  file_path: ''
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id' },
  { title: '名称', dataIndex: 'name', key: 'name' },
  { title: '版本', dataIndex: 'version', key: 'version' },
  { title: '作者', dataIndex: 'author', key: 'author' },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
  { title: '操作', key: 'action' }
]

const loadPlugins = async () => {
  loading.value = true
  try {
    const response = await pluginAPI.getAll()
    plugins.value = response.data
  } catch (error) {
    message.error('加载插件列表失败')
  } finally {
    loading.value = false
  }
}

const viewDetail = (id) => {
  router.push(`/plugins/${id}`)
}

const downloadPlugin = (plugin) => {
  message.success(`下载插件: ${plugin.name}`)
}

const deletePlugin = async (id) => {
  try {
    await pluginAPI.delete(id)
    message.success('删除成功')
    loadPlugins()
  } catch (error) {
    message.error('删除失败')
  }
}

const showModal = () => {
  visible.value = true
}

const handleOk = async () => {
  try {
    await pluginAPI.create(form.value)
    message.success('发布成功')
    visible.value = false
    form.value = {
      name: '',
      version: '',
      description: '',
      author: '',
      file_path: ''
    }
    loadPlugins()
  } catch (error) {
    message.error('发布失败')
  }
}

const handleCancel = () => {
  visible.value = false
}

onMounted(() => {
  loadPlugins()
})
</script>
