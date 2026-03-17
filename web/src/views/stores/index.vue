<template>
  <div class="stores">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>门店列表</span>
          <div class="filter">
            <el-select v-model="statusFilter" placeholder="状态筛选" clearable @change="handleFilter">
              <el-option label="待审批" value="pending" />
              <el-option label="已批准" value="approved" />
              <el-option label="已拒绝" value="rejected" />
            </el-select>
          </div>
        </div>
      </template>

      <el-table :data="storeStore.stores" v-loading="storeStore.loading">
        <el-table-column prop="name" label="门店名称" width="150" />
        <el-table-column prop="rustdesk_id" label="RustDesk ID" width="120" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="在线" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_online ? 'success' : 'info'" size="small">
              {{ row.is_online ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="CPU" width="80">
          <template #default="{ row }">
            {{ row.cpu_percent ? row.cpu_percent.toFixed(1) + '%' : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="内存" width="80">
          <template #default="{ row }">
            {{ row.memory_percent ? row.memory_percent.toFixed(1) + '%' : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="client_version" label="客户端版本" width="100" />
        <el-table-column prop="last_seen" label="最后在线" width="160">
          <template #default="{ row }">
            {{ row.last_seen ? formatTime(row.last_seen) : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="350">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'pending'"
              type="success"
              size="small"
              @click="handleApprove(row)"
            >
              审批
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="warning"
              size="small"
              @click="handleReject(row)"
            >
              拒绝
            </el-button>
            <el-button
              v-if="row.status === 'approved' && row.is_online"
              type="primary"
              size="small"
              @click="handleDesktop(row)"
            >
              远程桌面
            </el-button>
            <el-dropdown v-if="row.status === 'approved'" trigger="click">
              <el-button type="info" size="small">
                更多<el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleCommand(row)">命令执行</el-dropdown-item>
                  <el-dropdown-item @click="handleProcess(row)">进程监控</el-dropdown-item>
                  <el-dropdown-item @click="handleResource(row)">资源监控</el-dropdown-item>
                  <el-dropdown-item divided @click="handleShowToken(row)">查看TOKEN</el-dropdown-item>
                  <el-dropdown-item divided @click="handleDelete(row)" style="color: #f56c6c">删除门店</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="storeStore.total"
        layout="total, prev, pager, next"
        style="margin-top: 20px; justify-content: flex-end"
        @current-change="handlePageChange"
      />
    </el-card>

    <!-- 拒绝原因对话框 -->
    <el-dialog v-model="rejectDialogVisible" title="拒绝原因" width="400px">
      <el-input v-model="rejectReason" type="textarea" :rows="3" placeholder="请输入拒绝原因" />
      <template #footer>
        <el-button @click="rejectDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmReject">确认</el-button>
      </template>
    </el-dialog>

    <!-- TOKEN显示对话框 -->
    <el-dialog v-model="tokenDialogVisible" title="门店TOKEN" width="500px">
      <el-form label-width="100px">
        <el-form-item label="门店名称">
          <el-input :value="currentStore?.name" readonly />
        </el-form-item>
        <el-form-item label="TOKEN">
          <el-input :value="currentStore?.token" readonly>
            <template #append>
              <el-button @click="copyToken">复制</el-button>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="tokenDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="handleRegenerateToken">重新生成</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import { useStoreStore } from '@/stores/store'
import type { Store } from '@/api/store'

const router = useRouter()
const storeStore = useStoreStore()

const statusFilter = ref('')
const currentPage = ref(1)
const pageSize = 10

const rejectDialogVisible = ref(false)
const rejectReason = ref('')
const rejectStore = ref<Store | null>(null)

const tokenDialogVisible = ref(false)
const currentStore = ref<Store | null>(null)

function getStatusType(status: string) {
  const map: Record<string, string> = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return map[status] || 'info'
}

function getStatusText(status: string) {
  const map: Record<string, string> = {
    pending: '待审批',
    approved: '已批准',
    rejected: '已拒绝'
  }
  return map[status] || status
}

function formatTime(time: string) {
  return new Date(time).toLocaleString('zh-CN')
}

async function fetchData() {
  await storeStore.fetchStores({
    status: statusFilter.value || undefined,
    page: currentPage.value,
    page_size: pageSize
  })
}

function handleFilter() {
  currentPage.value = 1
  fetchData()
}

function handlePageChange(page: number) {
  currentPage.value = page
  fetchData()
}

async function handleApprove(store: Store) {
  await ElMessageBox.confirm(`确认审批门店 "${store.name}"？`, '审批确认')
  await storeStore.approveStore(store.id)
  ElMessage.success('审批成功')
  fetchData()
}

function handleReject(store: Store) {
  rejectStore.value = store
  rejectReason.value = ''
  rejectDialogVisible.value = true
}

async function confirmReject() {
  if (!rejectStore.value) return
  await storeStore.rejectStore(rejectStore.value.id, rejectReason.value)
  ElMessage.success('已拒绝')
  rejectDialogVisible.value = false
  fetchData()
}

async function handleDelete(store: Store) {
  await ElMessageBox.confirm(`确认删除门店 "${store.name}"？此操作不可恢复。`, '删除确认', {
    type: 'warning'
  })
  await storeStore.deleteStore(store.id)
  ElMessage.success('删除成功')
  fetchData()
}

function handleShowToken(store: Store) {
  currentStore.value = store
  tokenDialogVisible.value = true
}

function copyToken() {
  if (currentStore.value?.token) {
    navigator.clipboard.writeText(currentStore.value.token)
    ElMessage.success('已复制到剪贴板')
  }
}

async function handleRegenerateToken() {
  if (!currentStore.value) return
  await ElMessageBox.confirm('重新生成TOKEN后，原TOKEN将失效，确认继续？', '确认')
  const result = await storeStore.regenerateToken(currentStore.value.id)
  currentStore.value = { ...currentStore.value, token: result.token }
  ElMessage.success('TOKEN已重新生成')
}

function handleDesktop(store: Store) {
  router.push(`/stores/${store.id}/desktop`)
}

function handleCommand(store: Store) {
  router.push(`/stores/${store.id}/command`)
}

function handleProcess(store: Store) {
  router.push(`/stores/${store.id}/process`)
}

function handleResource(store: Store) {
  router.push(`/stores/${store.id}/resource`)
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.stores {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter {
  display: flex;
  gap: 10px;
}
</style>