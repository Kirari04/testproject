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

const name = ref('')
const crfFiles = ref<UploadFileInfo[]>([])
const keyFiles = ref<UploadFileInfo[]>([])

async function addCertificate() {
    if (crfFiles.value.length === 0 || keyFiles.value.length === 0) {
        useToast().error('Please select both certificate and key file')
        return
    }

    loadingBar.start()
    const formData = new FormData()
    formData.append('name', name.value)
    formData.append('crt', crfFiles.value[0].file!)
    formData.append('key', keyFiles.value[0].file!)
    await axios.post<string>(`${import.meta.env.VITE_APP_API}/api/certificate`, formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
        },
    })
        .then(() => {
            emit('onCreated')
            useToast().success('Certificate uploaded')
            showModal.value = false

            name.value = ''
            crfFiles.value = []
            keyFiles.value = []
        })
        .catch(err => {
            useToast().error(
                h(ToastDesc, {
                    title: 'Failed to upload certificate',
                    message: err.response.data ?? err.message,
                }), {
                timeout: 5000,
            })
        })
    loadingBar.finish()
}
</script>
<template>
    <n-button @click="showModal = true">
        Upload Certificate
    </n-button>
    <n-modal v-model:show="showModal">
        <n-card title="Upload Certificate" style="width: 700px;" :bordered="false" size="huge" role="dialog"
            aria-modal="true">
            <n-space vertical>
                <n-card>
                    <h3>Certificate name</h3>
                    <n-space>
                        <n-input v-model:value="name" type="text" placeholder="Certificate name" />
                    </n-space>
                    <h3>Certificate File (.crt)</h3>
                    <n-upload v-model:file-list="crfFiles" :max="1" :multiple="false">
                        <n-button>Select File</n-button>
                    </n-upload>
                    <h3>Key File (.key)</h3>
                    <n-upload v-model:file-list="keyFiles" :max="1" :multiple="false">
                        <n-button>Select File</n-button>
                    </n-upload>
                </n-card>
                <n-button type="primary" @click="addCertificate">Upload</n-button>
            </n-space>
        </n-card>
    </n-modal>
</template>