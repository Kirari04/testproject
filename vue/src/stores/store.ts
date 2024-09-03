import { ref, computed, defineComponent, h } from 'vue'
import { defineStore } from 'pinia'
import axios from 'axios'
import { useToast } from 'vue-toastification'
import ToastDesc from '@/components/ToastDesc.vue'


export const useStore = defineStore('store', () => {
  const isProxyRunning = ref(false)
  function setIsProxyRunning(v: boolean) {
    isProxyRunning.value = v
  }
  async function checkIsProxyRunning() {
    await axios.get<string>(`${import.meta.env.VITE_APP_API}/api/status`, {
      timeout: 1000 * 5,
    })
      .then(res => {
        isProxyRunning.value = (res.data === 'ok')
      })
      .catch(err => {
        isProxyRunning.value = false
        useToast().error(
          h(ToastDesc, {
            title: 'Proxy is not running',
            message: err.message,
          }), {
          timeout: 5000,
        })
      })
  }

  return {
    isProxyRunning,
    setIsProxyRunning,
    checkIsProxyRunning,
  }
})
