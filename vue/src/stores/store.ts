import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import axios from 'axios'

export const useStore = defineStore('store', () => {
  const isProxyRunning = ref(false)
  function setIsProxyRunning(v: boolean) {
    isProxyRunning.value = v
  }
  async function checkIsProxyRunning() {
    await axios.get<string>(`${import.meta.env.VITE_APP_API}/api/status`)
      .then(res => {
        isProxyRunning.value = (res.data === 'ok')
      })
      .catch(err => {
        isProxyRunning.value = false
        console.log(err.message)
      })
  }

  return {
    isProxyRunning,
    setIsProxyRunning,
    checkIsProxyRunning,
  }
})
