<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-value">{{ stats.totalStores }}</div>
            <div class="stat-label">门店总数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-value online">{{ stats.onlineStores }}</div>
            <div class="stat-label">在线门店</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-value pending">{{ stats.pendingStores }}</div>
            <div class="stat-label">待审批</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-value">{{ stats.offlineStores }}</div>
            <div class="stat-label">离线门店</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="recent-stores" style="margin-top: 20px">
      <template #header>
        <span>最近连接的门店</span>
      </template>
      <el-table :data="recentStores" v-loading="loading">
        <el-table-column prop="name" label="门店名称" />
        <el-table-column prop="rustdesk_id" label="RustDesk ID" />
        <el-table-column label="状态">
          <template #default="{ row }">
            <el-tag :type="row.is_online ? 'success' : 'info'">
              {{ row.is_online ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cpu_percent" label="CPU">
          <template #default="{ row }">
            {{ row.cpu_percent ? row.cpu_percent.toFixed(1) + '%' : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="memory_percent" label="内存">
          <template #default="{ row }">
            {{ row.memory_percent ? row.memory_percent.toFixed(1) + '%' : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="last_seen" label="最后在线">
          <template #default="{ row }">
            {{ row.last_seen ? formatTime(row.last_seen) : '-' }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useStoreStore } from '@/stores/store'
import type { Store } from '@/api/store'

const storeStore = useStoreStore()
const loading = ref(true)
const recentStores = ref<Store[]>([])

const stats = computed(() => {
  const total = storeStore.total
  const online = storeStore.stores.filter(s => s.is_online).length
  const pending = storeStore.stores.filter(s => s.status === 'pending').length
  return {
    totalStores: total,
    onlineStores: online,
    pendingStores: pending,
    offlineStores: total - online
  }
})

function formatTime(time: string) {
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(async () => {
  try {
    await storeStore.fetchStores({ page: 1, page_size: 10 })
    recentStores.value = storeStore.stores
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.stat-card {
  text-align: center;
  padding: 10px 0;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
}

.stat-value.online {
  color: #67c23a;
}

.stat-value.pending {
  color: #e6a23c;
}

.stat-label {
  margin-top: 8px;
  color: #909399;
}
</style>