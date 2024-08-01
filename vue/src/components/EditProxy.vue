<script setup lang="ts">
import { NCard, useLoadingBar, NSpace, NInput, NInputNumber, NButton, NIcon, NTabs, NTabPane, NSelect, NModal, NAlert, NSwitch, NFlex } from 'naive-ui'
import { onMounted, ref, h, watch } from 'vue'
import { type Component } from 'vue'
import {
    AddRound,
    MinusRound,
} from '@vicons/material'
import axios from 'axios';
import { useToast } from 'vue-toastification';
import ToastDesc from '@/components/ToastDesc.vue'
import type { Certificate, Frontend } from 'env';

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

const props = defineProps<{
    proxyId: number
}>()

const showModal = ref(false)
watch(showModal, () => {
    if (showModal.value) {
        getProxy()
    }
})

// const store = useStore()

const port = ref(80)
const https = ref(false)
const domain = ref('')
const bwInLimit = ref(0)
const bwInLimitUnit = ref(1 * 1024 * 1024)
const bwInPeriod = ref(1)
const bwOutLimit = ref(0)
const bwOutLimitUnit = ref(1 * 1024 * 1024)
const bwOutPeriod = ref(1)
const rateLimit = ref(0)
const ratePeriod = ref(1)
const hardRateLimit = ref(0)
const hardRatePeriod = ref(1)

const backendHttps = ref(false)
const backendHttpsVerify = ref(false)

const httpCheck = ref(false)
const httpCheckMethod = ref('GET')
const httpCheckPath = ref('/')
const httpCheckExpectStatus = ref(200)
const httpCheckInterval = ref(1)
const httpCheckFailAfter = ref(5)
const httpCheckRecoverAfter = ref(2)

const requestBodyLimit = ref(0)
const requestBodyLimitUnit = ref(1 * 1024 * 1024)

const backends = ref<{ addr: string }[]>([{ addr: '' }])
const aliases = ref<{ domain: string }[]>([{ domain: '' }])

watch(https, async () => {
    if (https.value && port.value === 80) {
        port.value = 443
    } else if(!https.value && port.value === 443) {
        port.value = 80
    }
})

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

function addAlias() {
    aliases.value.push({ domain: '' })
}
function removeAlias(index: number) {
    if (aliases.value.length === 1) return
    aliases.value.splice(index, 1)
}

function addBackend() {
    backends.value.push({ addr: '' })
}

function removeBackend(index: number) {
    if (backends.value.length === 1) return
    backends.value.splice(index, 1)
}

const isCreatingProxy = ref(false)
async function editProxy() {
    isCreatingProxy.value = true
    loadingBar.start()
    await axios.put<string>(`${import.meta.env.VITE_APP_API}/api/proxy`, {
        id: props.proxyId,
        port: port.value,
        https: https.value,
        domain: domain.value,
        bw_in_limit: bwInLimit.value,
        bw_in_limit_unit: bwInLimitUnit.value,
        bw_in_period: bwInPeriod.value,
        bw_out_limit: bwOutLimit.value,
        bw_out_limit_unit: bwOutLimitUnit.value,
        bw_out_period: bwOutPeriod.value,
        rate_limit: rateLimit.value,
        rate_period: ratePeriod.value,
        hard_rate_limit: hardRateLimit.value,
        hard_rate_period: hardRatePeriod.value,
        backend_https: backendHttps.value,
        backend_https_verify: backendHttpsVerify.value,
        http_check: httpCheck.value,
        http_check_method: httpCheckMethod.value,
        http_check_path: httpCheckPath.value,
        http_check_expect_status: httpCheckExpectStatus.value,
        http_check_interval: httpCheckInterval.value,
        http_check_fail_after: httpCheckFailAfter.value,
        http_check_recover_after: httpCheckRecoverAfter.value,
        request_body_limit: requestBodyLimit.value,
        request_body_limit_unit: requestBodyLimitUnit.value,
        backends: backends.value.map(b => ({ address: b.addr })).filter(b => b.address !== ''),
        aliases: aliases.value.map(a => ({ domain: a.domain })).filter(a => a.domain !== ''),
    })
        .then(() => {
            emit('onCreated')
            useToast().success('Proxy edited')
            showModal.value = false
        })
        .catch(err => {
            useToast().error(
                h(ToastDesc, {
                    title: 'Failed to edit proxy',
                    message: err.response.data ?? err.message,
                }), {
                timeout: 5000,
            })
        })
    loadingBar.finish()
    isCreatingProxy.value = false
}

const isLoadingProxy = ref(false)
async function getProxy() {
    loadingBar.start()
    isLoadingProxy.value = true
    await axios.get<Frontend>(`${import.meta.env.VITE_APP_API}/api/proxy`, {
        params: {
            id: props.proxyId,
        }
    })
        .then(data => {
            if(data.status === 200) {
                port.value = data.data.port
                https.value = data.data.https
                domain.value = data.data.domain
                bwInLimit.value = data.data.bw_limit
                bwInLimitUnit.value = data.data.bw_limit_unit
                bwInPeriod.value = data.data.bw_period
                bwOutLimit.value = data.data.bw_out_limit
                bwOutLimitUnit.value = data.data.bw_out_limit_unit
                bwOutPeriod.value = data.data.bw_out_period
                rateLimit.value = data.data.rate_limit
                ratePeriod.value = data.data.rate_period
                hardRateLimit.value = data.data.hard_rate_limit
                hardRatePeriod.value = data.data.hard_rate_period
                backendHttps.value = data.data.backends[0].https ?? false
                backendHttpsVerify.value = data.data.backends[0].https_verify ?? false
                httpCheck.value = data.data.http_check
                httpCheckMethod.value = data.data.http_check_method
                httpCheckPath.value = data.data.http_check_path
                httpCheckExpectStatus.value = data.data.http_check_expect_status
                httpCheckInterval.value = data.data.http_check_interval
                httpCheckFailAfter.value = data.data.http_check_fail_after
                httpCheckRecoverAfter.value = data.data.http_check_recover_after
                requestBodyLimit.value = data.data.request_body_limit
                requestBodyLimitUnit.value = data.data.request_body_limit_unit
                backends.value = data.data.backends.map(b => ({ addr: b.address }))
                aliases.value = data.data.aliases.map(a => ({ domain: a.domain }))
            }else{
                useToast().error(
                h(ToastDesc, {
                    title: 'Failed to load proxy',
                    message: data.statusText,
                }), {
                timeout: 5000,
            })
            }
        })
        .catch(err => {
            useToast().error(
                h(ToastDesc, {
                    title: 'Failed to load proxy',
                    message: err.response.data ?? err.message,
                }), {
                timeout: 5000,
            })
        })
    isLoadingProxy.value = false
    loadingBar.finish()
}
</script>
<template>
    <n-button @click="showModal = true" type="warning">
        Edit
    </n-button>
    <n-modal v-model:show="showModal">
        <n-card title="Edit Proxy" style="width: 700px;" :bordered="false" size="huge" role="dialog"
            aria-modal="true">
            <n-space vertical>
                <n-tabs type="line" animated>
                    <n-tab-pane name="details" tab="Details">
                        <n-space vertical>
                            <n-card>
                                <h3>Listen on :{{ port }}</h3>
                                <n-space vertical>
                                    <n-flex align="center" justify="start">
                                        Use HTTPS
                                        <n-switch v-model:value="https" />
                                    </n-flex>
                                    <n-space>
                                        <n-input-number v-model:value="port" placeholder="80" />
                                    </n-space>
                                </n-space>
                                <h3>Domain</h3>
                                <n-space>
                                    <n-input v-model:value="domain" type="text" placeholder="example.com" />
                                </n-space>
                                <h3>Aliases</h3>
                                <n-space vertical>
                                    <n-space v-for="(alias, i) in aliases" :key="i">
                                        <n-input v-model:value="alias.domain" type="text"
                                            placeholder="alias.example.com" />
                                        <n-button :render-icon="renderIcon(MinusRound)" :disabled="aliases.length === 1"
                                            @click="removeAlias(i)"></n-button>
                                    </n-space>
                                    <n-space>
                                        <n-button :render-icon="renderIcon(AddRound)" @click="addAlias()"></n-button>
                                    </n-space>
                                </n-space>
                            </n-card>
                            <n-card>
                                <h3>Backends</h3>
                                <n-space vertical>
                                    <n-space>
                                        Use HTTPS
                                        <n-switch v-model:value="backendHttps" />
                                        HTTPS Verify Certificate
                                        <n-switch v-model:value="backendHttpsVerify" :disabled="!backendHttps" />
                                    </n-space>
                                    <n-space v-for="(backend, i) in backends" :key="i">
                                        <n-input v-model:value="backend.addr" type="text"
                                            placeholder="127.0.0.1:8080" />
                                        <n-button :render-icon="renderIcon(MinusRound)"
                                            :disabled="backends.length === 1" @click="removeBackend(i)"></n-button>
                                    </n-space>
                                    <n-space>
                                        <n-button :render-icon="renderIcon(AddRound)" @click="addBackend()"></n-button>
                                    </n-space>
                                </n-space>
                            </n-card>
                        </n-space>
                    </n-tab-pane>
                    <n-tab-pane name="health-check" tab="Health check">
                        <n-space vertical>
                            <n-card>
                                <n-space vertical>
                                    <h3>Backend Health Check</h3>
                                    <n-space>
                                        Enable
                                        <n-switch v-model:value="httpCheck" />
                                    </n-space>
                                    <n-space>
                                        <div>
                                            Path
                                            <n-input v-model:value="httpCheckPath" :disabled="!httpCheck" type="text"
                                                placeholder="/" />
                                        </div>
                                    </n-space>
                                    <n-space>
                                        <div>
                                            Method
                                            <n-select v-model:value="httpCheckMethod" :disabled="!httpCheck"
                                                :options="['GET', 'POST', 'HEAD'].map(m => ({ label: m, value: m }))"
                                                style="min-width: 130px;" />
                                        </div>
                                    </n-space>
                                    <n-space>
                                        <div>
                                            Expect status
                                            <n-input-number v-model:value="httpCheckExpectStatus" :disabled="!httpCheck"
                                                placeholder="200" />
                                        </div>
                                    </n-space>
                                    <n-space>
                                        <div>
                                            Interval (seconds)
                                            <n-input-number v-model:value="httpCheckInterval" :disabled="!httpCheck"
                                                placeholder="1" />
                                        </div>
                                    </n-space>
                                    <n-space>
                                        <div>
                                            Fail after {{ httpCheckFailAfter }} requests
                                            <n-input-number v-model:value="httpCheckFailAfter" :disabled="!httpCheck"
                                                placeholder="5" />
                                        </div>
                                    </n-space>
                                    <n-space>
                                        <div>
                                            Recover after {{ httpCheckRecoverAfter }} requests
                                            <n-input-number v-model:value="httpCheckRecoverAfter" :disabled="!httpCheck"
                                                placeholder="2" />
                                        </div>
                                    </n-space>
                                </n-space>
                            </n-card>
                        </n-space>
                    </n-tab-pane>
                    <n-tab-pane name="bandwith" tab="Bandwith">
                        <n-space vertical>
                            <n-alert type="info">
                                If the limit is 0, the bandwidth limit won't be applied. <br>
                                The bandwith limit is applied on a per ip basis.
                            </n-alert>
                            <n-card>
                                <n-space vertical>
                                    <h3>Upload</h3>
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
                                    <h3>Download</h3>
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
                            </n-card>
                        </n-space>
                    </n-tab-pane>
                    <n-tab-pane name="rate-limit" tab="Rate limit">
                        <n-space vertical>
                            <n-alert type="info">
                                If the limit is 0, the Rate limit won't be applied. <br>
                                The rate limit is applied on a per ip basis. <br>
                                <strong>Soft limit</strong> Will respond with a http status code 429 if the limit is
                                reached. <br>
                                <strong>Hard limit</strong> Will silently drop the request if the limit is
                                reached.
                            </n-alert>
                            <n-card>
                                <h3>Soft limit</h3>
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
                                <h3>Hard limit</h3>
                                <n-space>
                                    <div>
                                        Limit
                                        <n-input-number v-model:value="hardRateLimit" placeholder="0" />
                                    </div>
                                    <div>
                                        Period (seconds)
                                        <n-input-number v-model:value="hardRatePeriod" placeholder="0" />
                                    </div>
                                </n-space>
                            </n-card>
                        </n-space>
                    </n-tab-pane>
                    <n-tab-pane name="other-limits" tab="Other Limits">
                        <n-space vertical>
                            <n-alert type="info">
                                If the limit is 0, the Body limit won't be applied. <br>
                            </n-alert>
                            <n-card>
                                <h3>Request Body Limit</h3>
                                <n-space>
                                    <div>
                                        Limit
                                        <n-input-number v-model:value="requestBodyLimit" placeholder="0" />
                                    </div>
                                    <div>
                                        Unit
                                        <n-select v-model:value="requestBodyLimitUnit"
                                            :options="bwUnits.map(u => ({ label: u.label, value: u.value }))"
                                            style="min-width: 130px;" />
                                    </div>
                                </n-space>
                            </n-card>
                        </n-space>
                    </n-tab-pane>
                </n-tabs>
                <n-button type="primary" @click="editProxy()" :loading="isCreatingProxy || isLoadingProxy">Save</n-button>
            </n-space>
        </n-card>
    </n-modal>
</template>