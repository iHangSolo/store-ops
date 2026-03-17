<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetConnectionStatus, GetStoreName, GetVersion, GetToken, SetToken, SetStoreName, GetDeviceID, GetRustDeskID } from '../../wailsjs/go/main/App'

const status = ref('unknown')
const storeName = ref('')
const version = ref('')
const token = ref('')
const deviceId = ref('')
const rustdeskId = ref('')
const showTokenDialog = ref(false)
const newToken = ref('')
const newStoreName = ref('')

onMounted(async () => {
  try {
    status.value = await GetConnectionStatus()
    storeName.value = await GetStoreName()
    version.value = await GetVersion()
    token.value = await GetToken()
    deviceId.value = await GetDeviceID()
    rustdeskId.value = await GetRustDeskID()
  } catch (e) {
    console.error('Failed to get app info:', e)
  }
})

async function saveToken() {
  if (newToken.value) {
    await SetToken(newToken.value)
    token.value = newToken.value
    showTokenDialog.value = false
    newToken.value = ''
    status.value = await GetConnectionStatus()
  }
}

async function saveStoreName() {
  if (newStoreName.value) {
    await SetStoreName(newStoreName.value)
    storeName.value = newStoreName.value
    newStoreName.value = ''
  }
}
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
        <span class="value" v-if="storeName">{{ storeName }}</span>
        <span class="value edit" v-else @click="newStoreName = ''; ">点击设置</span>
      </div>
      <div class="info-row">
        <span class="label">设备ID:</span>
        <span class="value small">{{ deviceId.substring(0, 8) }}...</span>
      </div>
      <div class="info-row">
        <span class="label">RustDesk:</span>
        <span class="value">{{ rustdeskId || '未获取' }}</span>
      </div>
      <div class="info-row">
        <span class="label">状态:</span>
        <span class="value" :class="status">{{ status === 'connected' ? '已连接' : '未连接' }}</span>
      </div>
      <div class="actions">
        <button class="btn" @click="showTokenDialog = true">设置TOKEN</button>
      </div>
    </div>

    <!-- Token对话框 -->
    <div class="dialog-overlay" v-if="showTokenDialog" @click.self="showTokenDialog = false">
      <div class="dialog">
        <h3>设置TOKEN</h3>
        <input v-model="newToken" type="text" placeholder="请输入TOKEN" />
        <div class="dialog-actions">
          <button class="btn" @click="showTokenDialog = false">取消</button>
          <button class="btn primary" @click="saveToken">保存</button>
        </div>
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
  min-width: 300px;
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

.info-row:last-of-type {
  border-bottom: none;
}

.label {
  color: rgba(255, 255, 255, 0.6);
}

.value {
  font-weight: 500;
}

.value.small {
  font-size: 12px;
}

.value.edit {
  color: #60a5fa;
  cursor: pointer;
}

.value.connected {
  color: #4ade80;
}

.value.disconnected {
  color: #f87171;
}

.actions {
  margin-top: 16px;
  text-align: center;
}

.btn {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: #fff;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.btn.primary {
  background: #3b82f6;
  border-color: #3b82f6;
}

.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.dialog {
  background: #1e293b;
  border-radius: 12px;
  padding: 20px;
  min-width: 280px;
}

.dialog h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
}

.dialog input {
  width: 100%;
  padding: 10px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  font-size: 14px;
  box-sizing: border-box;
}

.dialog input::placeholder {
  color: rgba(255, 255, 255, 0.4);
}

.dialog-actions {
  margin-top: 16px;
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}
</style>