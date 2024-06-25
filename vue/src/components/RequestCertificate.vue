<script setup lang="ts">
import { NCard, useLoadingBar, NSpace, NInput, NInputNumber, NButton, NIcon, NTabs, NTabPane, NSelect, NModal, NAlert, NUpload } from 'naive-ui'
import type { UploadFileInfo } from 'naive-ui'
import { onMounted, ref, h } from 'vue'
import { type Component } from 'vue'
import {
    AddRound,
    MinusRound,
} from '@vicons/material'
import axios from 'axios';
import { useToast } from 'vue-toastification';
import ToastDesc from '@/components/ToastDesc.vue'

function renderIcon(icon: Component) {
    return () => h(NIcon, null, { default: () => h(icon) })
}

const loadingBar = useLoadingBar()
onMounted(async () => {
    loadingBar.start()
    setTimeout(() => {
        loadingBar.finish()
    }, 500)
})

const emit = defineEmits<{
    onCreated: []
}>()

// const store = useStore()
const showModal = ref(false)
const isLoading = ref(false)
const domain = ref('')
async function reqCertificate() {
    isLoading.value = true
    loadingBar.start()
    await axios.post<string>(`${import.meta.env.VITE_APP_API}/api/certificate/request`, {
        domain: domain.value,
    })
        .then(() => {
            emit('onCreated')
            useToast().success('Certificate requested')
            showModal.value = false

            domain.value = ''
        })
        .catch(err => {
            useToast().error(
                h(ToastDesc, {
                    title: 'Failed to request certificate',
                    message: err.response.data ?? err.message,
                }), {
                timeout: 5000,
            })
        })
    loadingBar.finish()
    isLoading.value = false
}
</script>
<template>
    <n-button @click="showModal = true">
        Request Certificate
    </n-button>
    <n-modal v-model:show="showModal">
        <n-card title="Request Certificate" style="width: 700px;" :bordered="false" size="huge" role="dialog"
            aria-modal="true">
            <n-space vertical>
                <n-card>
                    <h3>Domain name</h3>
                    <n-space>
                        <n-input v-model:value="domain" type="text" placeholder="Domain name" />
                    </n-space>
                </n-card>
                <n-button type="primary" :loading="isLoading" @click="reqCertificate">Request</n-button>
            </n-space>
        </n-card>
    </n-modal>
</template>