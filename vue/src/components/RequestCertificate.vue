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
import type { Settings, SettingsAcmeCf } from 'env'

function renderIcon(icon: Component) {
    return () => h(NIcon, null, { default: () => h(icon) })
}

const loadingBar = useLoadingBar()
onMounted(async () => {
    loadingBar.start()
    await loadSettings()
    setTimeout(() => {
        loadingBar.finish()
    }, 500)
})

const emit = defineEmits<{
    onCreated: []
}>()
const acmeCloudflareDNSAPITokens = ref<SettingsAcmeCf[]>([])
async function loadSettings() {
    loadingBar.start()
    isLoading.value = true
    await axios.get<Settings>(`${import.meta.env.VITE_APP_API}/api/settings`)
        .then(res => {
            acmeCloudflareDNSAPITokens.value = res.data.acme.cf
        })
        .catch(err => {
            useToast().error(
                h(ToastDesc, {
                    title: 'Settings load failed',
                    message: err.response.data ?? err.message,
                }), {
                timeout: 5000,
            })
        })
    loadingBar.finish()
    isLoading.value = false
}

// const store = useStore()
const showModal = ref(false)
const isLoading = ref(false)
const domain = ref('')
async function reqCertificate() {
    isLoading.value = true
    loadingBar.start()
    await axios.post<string>(`${import.meta.env.VITE_APP_API}/api/certificate/request`, {
        domain: domain.value,
        auth_type: selectedAuthType.value,
        auth_id: selectedCredentialId.value,
    })
        .then(() => {
            emit('onCreated')
            useToast().success('Certificate obtained')
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

const options = [
    { label: 'Cloudflare DNS API', value: 'cloudflare_dns_api_token' },
]
const selectedAuthType = ref('cloudflare_dns_api_token')
const selectedCredentialId = ref<null | number>(null)


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
                    <n-space vertical>
                        <n-input v-model:value="domain" type="text" placeholder="Domain name" :disabled="isLoading" />
                        <n-select v-model:value="selectedAuthType" filterable placeholder="Select ACME Auth Type"
                            :options="options" :disabled="isLoading" />
                        <n-select v-model:value="selectedCredentialId" filterable placeholder="Select Credentials"
                            :options="acmeCloudflareDNSAPITokens.map(cred => ({ label: cred.name, value: cred.id }))"
                            :disabled="isLoading" />
                    </n-space>
                </n-card>
                <n-button type="primary" :loading="isLoading" @click="reqCertificate">Request</n-button>
            </n-space>
        </n-card>
    </n-modal>
</template>