<template>
  <div class="desktop">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ store?.name }} - 远程桌面</span>
          <el-button @click="router.back()">返回</el-button>
        </div>
      </template>

      <el-descriptions :column="3" border>
        <el-descriptions-item label="RustDesk ID">{{ store?.rustdesk_id || '-' }}</el-descriptions-item>
        <el-descriptions-item label="在线状态">
          <el-tag :type="store?.is_online ? 'success' : 'info'">
            {{ store?.is_online ? '在线' : '离线' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="当前状态">
          <el-tag v-if="status?.has_active_session" type="warning">被占用</el-tag>
          <el-tag v-else-if="status?.queue_count" type="info">排队中</el-tag>
          <el-tag v-else type="success">空闲</el-tag>
        </el-descriptions-item>
      </el-descriptions>

      <div class="actions" style="margin-top: 20px">
        <el-button
          v-if="store?.is_online"
          type="primary"
          :disabled="status?.has_active_session"
          @click="handleConnect"
        >
          连接远程桌面
        </el-button>
        <el-button
          v-if="status?.has_active_session"
          type="danger"
          @click="handleEndSession"
        >
          结束当前会话
        </el-button>
        <el-button
          v-if="store?.is_online && !status?.has_active_session && status?.queue_count"
          type="warning"
          @click="handleQueue"
        >
          加入排队 (当前{{ status?.queue_count }}人)
        </el-button>
      </div>

      <el-alert
        v-if="!store?.is_online"
        type="warning"
        title="门店离线"
        description="门店当前离线，无法进行远程桌面连接"
        show-icon
        style="margin-top: 20px"
      />

      <el-card v-if="status?.has_active_session && status?.current_session" style="margin-top: 20px">
        <template #header>
          <span>当前会话</span>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="操作人">{{ status?.current_session?.operator }}</el-descriptions-item>
          <el-descriptions-item label="开始时间">
            {{ formatTime(status?.current_session?.started_at) }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useStoreStore } from '@/stores/store'
import { getDesktopStatus, connectDesktop, endSession } from '@/api/desktop'
import type { DesktopStatus } from '@/api/desktop'

const route = useRoute()
const router = useRouter()
const storeStore = useStoreStore()

const storeId = route.params.id as string
const store = ref(storeStore.currentStore)
const status = ref<DesktopStatus | null>(null)

function formatTime(time: string) {
  return new Date(time).toLocaleString('zh-CN')
}

async function fetchStatus() {
  status.value = await getDesktopStatus(storeId)
}

async function handleConnect() {
  const result = await connectDesktop(storeId)
  window.open(result.rustdesk_url, '_blank')
  ElMessage.success('已打开 RustDesk 连接')
}

async function handleEndSession() {
  if (!status.value?.current_session) return
  await endSession(storeId, status.value.current_session.id)
  ElMessage.success('会话已结束')
  fetchStatus()
}

async function handleQueue() {
  ElMessage.info('排队功能开发中')
}

onMounted(async () => {
  if (!store.value) {
    await storeStore.fetchStore(storeId)
    store.value = storeStore.currentStore
  }
  fetchStatus()
})
</script>

<style scoped>
.desktop {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.actions {
  display: flex;
  gap: 10px;
}
</style>