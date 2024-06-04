<script setup lang="ts">
import { NCard, useLoadingBar, NSpace, NInput, NInputNumber, NButton, NIcon, NTabs, NTabPane, NSelect, NModal } from 'naive-ui'
import { onMounted, ref, h, defineEmits } from 'vue'
import { type Component } from 'vue'
import {
    AddRound,
    MinusRound,
} from '@vicons/material'
import axios from 'axios';
import { useToast } from 'vue-toastification';
import { useStore } from '@/stores/store';
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

const port = ref(80)
const domain = ref('')
const bwInLimit = ref(0)
const bwInLimitUnit = ref(1 * 1024 * 1024)
const bwInPeriod = ref(1)
const bwOutLimit = ref(0)
const bwOutLimitUnit = ref(1 * 1024 * 1024)
const bwOutPeriod = ref(1)
const rateLimit = ref(0)
const ratePeriod = ref(1)
const backends = ref<{ addr: string }[]>([{ addr: '' }])

const bwUnits = [{
    label: 'Bytes',
    value: 1,
}, {
    label: 'Kilobytes',
    value: 1 * 1024,
}, {
    label: 'Megabytes',
    value: 1 * 1024 * 1024,
}, {
    label: 'Gigabytes',
    value: 1 * 1024 * 1024 * 1024,
}]

function addBackend() {
    backends.value.push({ addr: '' })
}

function removeBackend(index: number) {
    if (backends.value.length === 1) return
    backends.value.splice(index, 1)
}

async function createProxy() {
    loadingBar.start()
    await axios.post<string>(`${import.meta.env.VITE_APP_API}/api/proxy`, {
        port: port.value,
        domain: domain.value,
        bw_in_limit: bwInLimit.value,
        bw_in_limit_unit: bwInLimitUnit.value,
        bw_in_period: bwInPeriod.value,
        bw_out_limit: bwOutLimit.value,
        bw_out_limit_unit: bwOutLimitUnit.value,
        bw_out_period: bwOutPeriod.value,
        rate_limit: rateLimit.value,
        rate_period: ratePeriod.value,
        backends: backends.value.map(b => ({ address: b.addr })),
    })
        .then(() => {
            emit('onCreated')
            useToast().success('Proxy created')
            showModal.value = false

            port.value = 80
            domain.value = ''
            bwInLimit.value = 0
            bwInLimitUnit.value = 1 * 1024 * 1024
            bwInPeriod.value = 1
            bwOutLimit.value = 0
            bwOutLimitUnit.value = 1 * 1024 * 1024
            bwOutPeriod.value = 1
            rateLimit.value = 0
            ratePeriod.value = 1
            backends.value = [{ addr: '' }]
        })
        .catch(err => {
            useToast().error(
                h(ToastDesc, {
                    title: 'Failed to create proxy',
                    message: err.message,
                }), {
                timeout: 5000,
            })
        })
    loadingBar.finish()
}
</script>
<template>
    <n-button @click="showModal = true">
        Create Proxy
    </n-button>
    <n-modal v-model:show="showModal">
        <n-card title="Create Proxy" style="width: 600px" :bordered="false" size="huge" role="dialog" aria-modal="true">
            <n-space vertical>
                <n-tabs type="line" animated>
                    <n-tab-pane name="details" tab="Details">
                        <n-space vertical>
                            <strong>Listen on :{{ port }}</strong>
                            <n-space>
                                <n-input-number v-model:value="port" placeholder="80" />
                            </n-space>
                            <strong>Domain</strong>
                            <n-space>
                                <n-input v-model:value="domain" type="text" placeholder="example.com" />
                            </n-space>
                            <strong>Backends</strong>
                            <n-space vertical>
                                <n-space v-for="(backend, i) in backends" :key="i">
                                    <n-input v-model:value="backend.addr" type="text" placeholder="127.0.0.1:8080" />
                                    <n-button :render-icon="renderIcon(MinusRound)" :disabled="backends.length === 1"
                                        @click="removeBackend(i)"></n-button>
                                </n-space>
                                <n-space>
                                    <n-button :render-icon="renderIcon(AddRound)" @click="addBackend()"></n-button>
                                </n-space>
                            </n-space>
                        </n-space>
                    </n-tab-pane>
                    <n-tab-pane name="bandwith" tab="Bandwith">
                        <n-space vertical>
                            <div>If the limit is 0, the bandwidth is unlimited.</div>
                            <n-space vertical>
                                <strong>Upload</strong>
                                <n-space>
                                    <div>
                                        Limit
                                        <n-input-number v-model:value="bwInLimit" placeholder="0" />
                                    </div>
                                    <div>
                                        Unit
                                        <n-select v-model:value="bwInLimitUnit"
                                            :options="bwUnits.map(u => ({ label: u.label, value: u.value }))"
                                            style="min-width: 130px;" />
                                    </div>
                                    <div>
                                        Period (seconds)
                                        <n-input-number v-model:value="bwInPeriod" placeholder="0" />
                                    </div>
                                </n-space>
                                <strong>Download</strong>
                                <n-space>
                                    <div>
                                        Limit
                                        <n-input-number v-model:value="bwOutLimit" placeholder="0" />
                                    </div>
                                    <div>
                                        Unit
                                        <n-select v-model:value="bwOutLimitUnit"
                                            :options="bwUnits.map(u => ({ label: u.label, value: u.value }))"
                                            style="min-width: 130px;" />
                                    </div>
                                    <div>
                                        Period (seconds)
                                        <n-input-number v-model:value="bwOutPeriod" placeholder="0" />
                                    </div>
                                </n-space>
                            </n-space>

                        </n-space>
                    </n-tab-pane>
                    <n-tab-pane name="rate-limit" tab="Rate limit">
                        <n-space vertical>
                            <div>If the limit is 0, the Rate limit is unlimited.</div>
                            <n-space>
                                <div>
                                    Limit
                                    <n-input-number v-model:value="rateLimit" placeholder="0" />
                                </div>
                                <div>
                                    Period (seconds)
                                    <n-input-number v-model:value="ratePeriod" placeholder="0" />
                                </div>
                            </n-space>
                        </n-space>
                    </n-tab-pane>
                </n-tabs>
                <n-button type="primary" @click="createProxy()">Save</n-button>
            </n-space>
        </n-card>
    </n-modal>
</template>