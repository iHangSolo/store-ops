import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as storeApi from '@/api/store'
import type { Store, StoreListParams } from '@/api/store'

export const useStoreStore = defineStore('store', () => {
  const stores = ref<Store[]>([])
  const total = ref(0)
  const loading = ref(false)
  const currentStore = ref<Store | null>(null)

  async function fetchStores(params: StoreListParams) {
    loading.value = true
    try {
      const result = await storeApi.getStores(params)
      stores.value = result.list
      total.value = result.total
    } finally {
      loading.value = false
    }
  }

  async function fetchStore(id: string) {
    loading.value = true
    try {
      currentStore.value = await storeApi.getStore(id)
    } finally {
      loading.value = false
    }
  }

  async function approveStore(id: string) {
    await storeApi.approveStore(id)
  }

  async function rejectStore(id: string, reason: string) {
    await storeApi.rejectStore(id, reason)
  }

  async function deleteStore(id: string) {
    await storeApi.deleteStore(id)
  }

  async function regenerateToken(id: string) {
    return await storeApi.regenerateToken(id)
  }

  return {
    stores,
    total,
    loading,
    currentStore,
    fetchStores,
    fetchStore,
    approveStore,
    rejectStore,
    deleteStore,
    regenerateToken
  }
})