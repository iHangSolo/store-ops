<template>
  <div class="command">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ store?.name }} - 命令执行</span>
          <el-button @click="router.back()">返回</el-button>
        </div>
      </template>

      <el-alert
        v-if="!store?.is_online"
        type="warning"
        title="门店离线"
        description="门店当前离线，无法执行命令"
        show-icon
        style="margin-bottom: 20px"
      />

      <el-form :model="form" label-width="80px">
        <el-form-item label="命令">
          <el-input
            v-model="form.command"
            type="textarea"
            :rows="3"
            placeholder="请输入要执行的命令"
            :disabled="!store?.is_online || executing"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            :loading="executing"
            :disabled="!store?.is_online || !form.command"
            @click="handleExecute"
          >
            执行
          </el-button>
          <el-button @click="form.command = ''">清空</el-button>
        </el-form-item>
      </el-form>

      <el-divider />

      <div v-if="result" class="result">
        <h4>执行结果</h4>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="状态">
            <el-tag :type="getResultStatusType(result.status)">
              {{ result.status }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="耗时">{{ result.duration_ms }}ms</el-descriptions-item>
          <el-descriptions-item label="执行时间">
            {{ formatTime(result.created_at) }}
          </el-descriptions-item>
        </el-descriptions>
        <el-input
          :model-value="result.output"
          type="textarea"
          :rows="10"
          readonly
          style="margin-top: 10px"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useStoreStore } from '@/stores/store'
import { executeCommand, getCommandResult } from '@/api/command'
import type { CommandResult } from '@/api/command'

const route = useRoute()
const router = useRouter()
const storeStore = useStoreStore()

const storeId = route.params.id as string
const store = ref(storeStore.currentStore)
const executing = ref(false)
const result = ref<CommandResult | null>(null)

const form = reactive({
  command: ''
})

function getResultStatusType(status: string) {
  const map: Record<string, string> = {
    success: 'success',
    failed: 'danger',
    timeout: 'warning',
    running: 'primary',
    pending: 'info'
  }
  return map[status] || 'info'
}

function formatTime(time: string) {
  return new Date(time).toLocaleString('zh-CN')
}

async function handleExecute() {
  if (!form.command) return

  executing.value = true
  result.value = null

  try {
    const { command_id } = await executeCommand(storeId, { command: form.command })
    ElMessage.success('命令已发送，等待结果...')

    // 轮询获取结果
    let attempts = 0
    const maxAttempts = 30
    const poll = async () => {
      const res = await getCommandResult(storeId, command_id)
      if (res.status !== 'pending' && res.status !== 'running') {
        result.value = res
        executing.value = false
        return
      }
      attempts++
      if (attempts < maxAttempts) {
        setTimeout(poll, 1000)
      } else {
        result.value = res
        executing.value = false
        ElMessage.warning('命令执行超时')
      }
    }
    poll()
  } catch {
    executing.value = false
  }
}

onMounted(async () => {
  if (!store.value) {
    await storeStore.fetchStore(storeId)
    store.value = storeStore.currentStore
  }
})
</script>

<style scoped>
.command {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.result h4 {
  margin-bottom: 10px;
  color: #303133;
}
</style>