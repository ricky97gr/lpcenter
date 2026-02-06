<template>
  <div>
    <a-page-header title="申请授权" />
    <a-card>
      <a-form :model="form" layout="vertical" @finish="handleSubmit">
        <a-form-item
          label="用户ID"
          name="user_id"
          :rules="[{ required: true, message: '请输入用户ID' }]"
        >
          <a-input v-model:value="form.user_id" placeholder="请输入用户ID" />
        </a-form-item>
        <a-form-item
          label="插件ID"
          name="plugin_id"
          :rules="[{ required: true, message: '请输入插件ID' }]"
        >
          <a-input v-model:value="form.plugin_id" placeholder="请输入插件ID" />
        </a-form-item>
        <a-form-item
          label="申请理由"
          name="reason"
          :rules="[{ required: true, message: '请输入申请理由' }]"
        >
          <a-textarea
            v-model:value="form.reason"
            placeholder="请输入申请理由"
            :rows="4"
          />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading">
            提交申请
          </a-button>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { licenseAPI } from '../api'

const loading = ref(false)
const form = ref({
  user_id: '',
  plugin_id: '',
  reason: ''
})

const handleSubmit = async () => {
  loading.value = true
  try {
    await licenseAPI.create({
      user_id: parseInt(form.value.user_id),
      plugin_id: parseInt(form.value.plugin_id),
      reason: form.value.reason
    })
    message.success('申请提交成功')
    form.value = {
      user_id: '',
      plugin_id: '',
      reason: ''
    }
  } catch (error) {
    message.error('申请提交失败')
  } finally {
    loading.value = false
  }
}
</script>
