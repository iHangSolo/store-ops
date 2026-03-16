<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetConnectionStatus, GetStoreName, GetVersion } from '../../wailsjs/go/main/App'

const status = ref('unknown')
const storeName = ref('')
const version = ref('')

onMounted(async () => {
  try {
    status.value = await GetConnectionStatus()
    storeName.value = await GetStoreName()
    version.value = await GetVersion()
  } catch (e) {
    console.error('Failed to get app info:', e)
  }
})
</script>

<template>
  <div class="container">
    <div class="status-card">
      <h2>门店运维客户端</h2>
      <div class="info-row">
        <span class="label">版本:</span>
        <span class="value">{{ version }}</span>
      </div>
      <div class="info-row">
        <span class="label">门店:</span>
        <span class="value">{{ storeName || '未设置' }}</span>
      </div>
      <div class="info-row">
        <span class="label">状态:</span>
        <span class="value" :class="status">{{ status === 'connected' ? '已连接' : '未连接' }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.container {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  color: #fff;
}

.status-card {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 24px;
  min-width: 280px;
  backdrop-filter: blur(10px);
}

h2 {
  margin: 0 0 16px 0;
  font-size: 18px;
  text-align: center;
}

.info-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.info-row:last-child {
  border-bottom: none;
}

.label {
  color: rgba(255, 255, 255, 0.6);
}

.value {
  font-weight: 500;
}

.value.connected {
  color: #4ade80;
}

.value.disconnected {
  color: #f87171;
}
</style>