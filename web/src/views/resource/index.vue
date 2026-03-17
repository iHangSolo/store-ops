<template>
  <div class="resource">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ store?.name }} - 资源监控</span>
          <div>
            <el-button @click="fetchResource" :loading="loading">刷新</el-button>
            <el-button @click="router.back()">返回</el-button>
          </div>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-card shadow="hover">
            <div class="stat-card">
              <div class="stat-label">CPU 使用率</div>
              <div class="stat-value" :class="{ warning: resource?.cpu_percent > 80 }">
                {{ resource?.cpu_percent?.toFixed(1) || 0 }}%
              </div>
              <el-progress
                :percentage="resource?.cpu_percent || 0"
                :stroke-width="10"
                :color="getProgressColor(resource?.cpu_percent)"
              />
            </div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card shadow="hover">
            <div class="stat-card">
              <div class="stat-label">内存使用率</div>
              <div class="stat-value" :class="{ warning: resource?.memory_percent > 80 }">
                {{ resource?.memory_percent?.toFixed(1) || 0 }}%
              </div>
              <el-progress
                :percentage="resource?.memory_percent || 0"
                :stroke-width="10"
                :color="getProgressColor(resource?.memory_percent)"
              />
              <div class="stat-sub">可用: {{ resource?.memory_available_gb?.toFixed(1) || 0 }} GB</div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-card style="margin-top: 20px">
        <template #header>
          <span>磁盘使用情况</span>
        </template>
        <el-table :data="resource?.disks || []">
          <el-table-column prop="name" label="磁盘" width="100" />
          <el-table-column label="使用率" width="200">
            <template #default="{ row }">
              <el-progress
                :percentage="row.percent"
                :stroke-width="10"
                :color="getProgressColor(row.percent)"
              />
            </template>
          </el-table-column>
          <el-table-column label="总容量">
            <template #default="{ row }">
              {{ row.total_gb?.toFixed(1) }} GB
            </template>
          </el-table-column>
          <el-table-column label="可用空间">
            <template #default="{ row }">
              {{ row.available_gb?.toFixed(1) }} GB
            </template>
          </el-table-column>
        </el-table>
      </el-card>

      <el-card style="margin-top: 20px">
        <template #header>
          <div style="display: flex; justify-content: space-between; align-items: center">
            <span>历史数据</span>
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              @change="fetchHistory"
            />
          </div>
        </template>
        <el-table :data="history" v-loading="historyLoading" max-height="400">
          <el-table-column prop="recorded_at" label="时间" width="180">
            <template #default="{ row }">
              {{ formatTime(row.recorded_at) }}
            </template>
          </el-table-column>
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
          <el-table-column prop="memory_available_gb" label="可用内存">
            <template #default="{ row }">
              {{ row.memory_available_gb?.toFixed(1) }} GB
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useStoreStore } from '@/stores/store'
import { getResource, getResourceHistory } from '@/api/resource'
import type { ResourceInfo } from '@/api/resource'

const route = useRoute()
const router = useRouter()
const storeStore = useStoreStore()

const storeId = route.params.id as string
const store = ref(storeStore.currentStore)
const loading = ref(false)
const historyLoading = ref(false)
const resource = ref<ResourceInfo | null>(null)
const history = ref<ResourceInfo[]>([])
const dateRange = ref<string[]>([])

function getProgressColor(percent: number | undefined) {
  if (!percent) return '#67c23a'
  if (percent > 80) return '#f56c6c'
  if (percent > 60) return '#e6a23c'
  return '#67c23a'
}

function formatTime(time: string) {
  return new Date(time).toLocaleString('zh-CN')
}

async function fetchResource() {
  loading.value = true
  try {
    resource.value = await getResource(storeId)
  } finally {
    loading.value = false
  }
}

async function fetchHistory() {
  historyLoading.value = true
  try {
    const params: Record<string, unknown> = { page: 1, page_size: 100 }
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_time = dateRange.value[0]
      params.end_time = dateRange.value[1]
    }
    const result = await getResourceHistory(storeId, params)
    history.value = result.list
  } finally {
    historyLoading.value = false
  }
}

onMounted(async () => {
  if (!store.value) {
    await storeStore.fetchStore(storeId)
    store.value = storeStore.currentStore
  }
  fetchResource()
  fetchHistory()
})
</script>

<style scoped>
.resource {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-card {
  text-align: center;
  padding: 20px 0;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 10px;
}

.stat-value {
  font-size: 36px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 10px;
}

.stat-value.warning {
  color: #f56c6c;
}

.stat-sub {
  font-size: 12px;
  color: #909399;
  margin-top: 10px;
}
</style>