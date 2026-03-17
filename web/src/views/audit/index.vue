<template>
  <div class="audit-logs">
    <el-card>
      <template #header>
        <span>审计日志</span>
      </template>

      <el-form :inline="true" class="filter-form">
        <el-form-item label="门店">
          <el-select v-model="filters.store_id" placeholder="选择门店" clearable>
            <el-option
              v-for="store in stores"
              :key="store.id"
              :label="store.name"
              :value="store.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="操作类型">
          <el-select v-model="filters.action_type" placeholder="选择类型" clearable>
            <el-option label="客户端连接" value="client_connect" />
            <el-option label="远程桌面连接" value="desktop_connect" />
            <el-option label="命令执行" value="command_execute" />
            <el-option label="软件推送" value="software_push" />
            <el-option label="终止进程" value="kill_process" />
            <el-option label="服务控制" value="service_control" />
            <el-option label="门店审批" value="store_approve" />
            <el-option label="门店拒绝" value="store_reject" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">查询</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="logs" v-loading="loading">
        <el-table-column prop="created_at" label="时间" width="160">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="store_name" label="门店" width="150">
          <template #default="{ row }">
            {{ row.store_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="action" label="操作" width="120" />
        <el-table-column prop="operator" label="操作人" width="100" />
        <el-table-column prop="detail" label="详情" show-overflow-tooltip />
      </el-table>

      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        style="margin-top: 20px; justify-content: flex-end"
        @current-change="handlePageChange"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { getAuditLogs } from '@/api/audit'
import { useStoreStore } from '@/stores/store'
import type { AuditLog } from '@/api/audit'

const storeStore = useStoreStore()

const logs = ref<AuditLog[]>([])
const total = ref(0)
const loading = ref(false)
const currentPage = ref(1)
const pageSize = 20

const filters = ref({
  store_id: '',
  action_type: ''
})
const dateRange = ref<string[]>([])

const stores = computed(() => storeStore.stores)

function formatTime(time: string) {
  return new Date(time).toLocaleString('zh-CN')
}

async function fetchData() {
  loading.value = true
  try {
    const params: Record<string, unknown> = {
      page: currentPage.value,
      page_size: pageSize
    }
    if (filters.value.store_id) {
      params.store_id = filters.value.store_id
    }
    if (filters.value.action_type) {
      params.action_type = filters.value.action_type
    }
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }

    const result = await getAuditLogs(params)
    logs.value = result.list
    total.value = result.total
  } finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  currentPage.value = page
  fetchData()
}

onMounted(() => {
  storeStore.fetchStores({ page: 1, page_size: 100 })
  fetchData()
})
</script>

<style scoped>
.audit-logs {
  padding: 20px;
}

.filter-form {
  margin-bottom: 20px;
}
</style>