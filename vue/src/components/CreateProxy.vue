<script setup lang="ts">
import { NCard, useLoadingBar, NSpace, NInput, NInputNumber, NButton, NIcon, NTabs, NTabPane, NSelect, NModal, NAlert, NSwitch } from 'naive-ui'
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
const hardRateLimit = ref(0)
const hardRatePeriod = ref(1)

const https = ref(false)
const httpsVerify = ref(false)

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
        hard_rate_limit: hardRateLimit.value,
        hard_rate_period: hardRatePeriod.value,
        https: https.value,
        https_verify: httpsVerify.value,
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
            https.value = false
            httpsVerify.value = false
            httpCheck.value = false
            httpCheckMethod.value = 'GET'
            httpCheckPath.value = '/'
            httpCheckExpectStatus.value = 200
            httpCheckInterval.value = 1
            httpCheckFailAfter.value = 5
            httpCheckRecoverAfter.value = 2
            requestBodyLimit.value = 0
            requestBodyLimitUnit.value = 1 * 1024 * 1024
            backends.value = [{ addr: '' }]
            aliases.value = [{ domain: '' }]
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
        <n-card title="Create Proxy" style="width: 700px;" :bordered="false" size="huge" role="dialog"
            aria-modal="true">
            <n-space vertical>
                <n-tabs type="line" animated>
                    <n-tab-pane name="details" tab="Details">
                        <n-space vertical>
                            <n-card>
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
                                <h3>Listen on :{{ port }}</h3>
                                <n-space>
                                    <n-input-number v-model:value="port" placeholder="80" />
                                </n-space>
                                <h3>Backends</h3>
                                <n-space vertical>
                                    <n-space>
                                        Use HTTPS
                                        <n-switch v-model:value="https" />
                                        HTTPS Verify Certificate
                                        <n-switch v-model:value="httpsVerify" :disabled="!https" />
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
                                            Fail after {{httpCheckFailAfter}} requests
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
                <n-button type="primary" @click="createProxy()">Save</n-button>
            </n-space>
        </n-card>
    </n-modal>
</template>