<template>
  <div class="process">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ store?.name }} - 进程监控</span>
          <div>
            <el-button @click="fetchProcesses" :loading="loading">刷新</el-button>
            <el-button @click="router.back()">返回</el-button>
          </div>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="进程列表" name="processes">
          <el-table :data="processes" v-loading="loading" max-height="500">
            <el-table-column prop="pid" label="PID" width="100" />
            <el-table-column prop="name" label="进程名" show-overflow-tooltip />
            <el-table-column prop="cpu_percent" label="CPU" width="100">
              <template #default="{ row }">
                {{ row.cpu_percent?.toFixed(1) }}%
              </template>
            </el-table-column>
            <el-table-column prop="memory_percent" label="内存" width="100">
              <template #default="{ row }">
                {{ row.memory_percent?.toFixed(1) }}%
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100" />
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button type="danger" size="small" @click="handleKillProcess(row.pid)">
                  终止
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="服务列表" name="services">
          <el-table :data="services" v-loading="loading" max-height="500">
            <el-table-column prop="name" label="服务名" width="200" show-overflow-tooltip />
            <el-table-column prop="display_name" label="显示名" show-overflow-tooltip />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'running' ? 'success' : 'info'">
                  {{ row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="start_type" label="启动类型" width="100" />
            <el-table-column label="操作" width="160">
              <template #default="{ row }">
                <el-button
                  v-if="row.status !== 'running'"
                  type="success"
                  size="small"
                  @click="handleStartService(row.name)"
                >
                  启动
                </el-button>
                <el-button
                  v-if="row.status === 'running'"
                  type="warning"
                  size="small"
                  @click="handleStopService(row.name)"
                >
                  停止
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useStoreStore } from '@/stores/store'
import { getProcesses, killProcess, getServices, startService, stopService } from '@/api/process'
import type { Process, Service } from '@/api/process'

const route = useRoute()
const router = useRouter()
const storeStore = useStoreStore()

const storeId = route.params.id as string
const store = ref(storeStore.currentStore)
const loading = ref(false)
const activeTab = ref('processes')
const processes = ref<Process[]>([])
const services = ref<Service[]>([])

async function fetchProcesses() {
  loading.value = true
  try {
    processes.value = await getProcesses(storeId)
  } finally {
    loading.value = false
  }
}

async function fetchServices() {
  loading.value = true
  try {
    services.value = await getServices(storeId)
  } finally {
    loading.value = false
  }
}

async function handleKillProcess(pid: number) {
  await ElMessageBox.confirm(`确认终止进程 PID: ${pid}？`, '确认终止', { type: 'warning' })
  await killProcess(storeId, pid)
  ElMessage.success('已发送终止命令')
  fetchProcesses()
}

async function handleStartService(name: string) {
  await startService(storeId, name)
  ElMessage.success('已发送启动命令')
  fetchServices()
}

async function handleStopService(name: string) {
  await ElMessageBox.confirm(`确认停止服务: ${name}？`, '确认停止', { type: 'warning' })
  await stopService(storeId, name)
  ElMessage.success('已发送停止命令')
  fetchServices()
}

onMounted(async () => {
  if (!store.value) {
    await storeStore.fetchStore(storeId)
    store.value = storeStore.currentStore
  }
  fetchProcesses()
})
</script>

<style scoped>
.process {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>