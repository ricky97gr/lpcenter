<template>
  <div>
    <a-page-header
      title="插件详情"
      @back="() => $router.push('/plugins')"
    />
    <a-spin :spinning="loading">
      <a-descriptions v-if="plugin" bordered :column="2">
        <a-descriptions-item label="ID">{{ plugin.id }}</a-descriptions-item>
        <a-descriptions-item label="名称">{{ plugin.name }}</a-descriptions-item>
        <a-descriptions-item label="版本">{{ plugin.version }}</a-descriptions-item>
        <a-descriptions-item label="作者">{{ plugin.author }}</a-descriptions-item>
        <a-descriptions-item label="描述" :span="2">
          {{ plugin.description }}
        </a-descriptions-item>
        <a-descriptions-item label="文件路径">
          {{ plugin.file_path }}
        </a-descriptions-item>
        <a-descriptions-item label="创建时间">
          {{ plugin.created_at }}
        </a-descriptions-item>
      </a-descriptions>
      <div style="margin-top: 16px">
        <a-button type="primary" @click="downloadPlugin">下载插件</a-button>
        <a-button @click="applyLicense" style="margin-left: 8px">
          申请授权
        </a-button>
      </div>
    </a-spin>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { pluginAPI } from '../api'

const route = useRoute()
const router = useRouter()
const plugin = ref(null)
const loading = ref(false)

const loadPlugin = async () => {
  loading.value = true
  try {
    const response = await pluginAPI.getById(route.params.id)
    plugin.value = response.data
  } catch (error) {
    message.error('加载插件详情失败')
  } finally {
    loading.value = false
  }
}

const downloadPlugin = () => {
  message.success(`下载插件: ${plugin.value.name}`)
}

const applyLicense = () => {
  router.push('/license-request')
}

onMounted(() => {
  loadPlugin()
})
</script>
