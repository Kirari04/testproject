<script setup lang="ts">
import axios from 'axios';
import { NCard, useLoadingBar, NTag, NH3, NH4, NFlex, NButton, NSpace, NAlert, NCode } from 'naive-ui'
import { onMounted, ref, h } from 'vue'
import { useToast } from 'vue-toastification';
import { useStore } from '@/stores/store'
import ToastDesc from '@/components/ToastDesc.vue'
import type { HaproxyCrashReasonsData } from 'env';

const loadingBar = useLoadingBar()
onMounted(async () => {
	loadingBar.start()
	await getCrashReport()
	setTimeout(() => {
		loadingBar.finish()
	}, 500)
})

const store = useStore()

const isProxyStarting = ref(false)
async function startProxy() {
	isProxyStarting.value = true
	loadingBar.start()
	await axios.get<string>(`${import.meta.env.VITE_APP_API}/api/start`).then(() => {
		useToast().success('Proxy started')
	}).catch(err => {
		useToast().error(
			h(ToastDesc, {
				title: 'Proxy start failed',
				message: err.response.data ?? err.message,
			}), {
			timeout: 5000,
		})
	})
	await store.checkIsProxyRunning()
	await getCrashReport()
	loadingBar.finish()
	isProxyStarting.value = false
}

const isProxyReloading = ref(false)
async function reloadProxy() {
	isProxyReloading.value = true
	loadingBar.start()
	await axios.get<string>(`${import.meta.env.VITE_APP_API}/api/reload`).then(() => {
		useToast().success('Proxy reloaded')
	}).catch(err => {
		useToast().error(
			h(ToastDesc, {
				title: 'Proxy reload failed',
				message: err.response.data ?? err.message,
			}), {
			timeout: 5000,
		})
	})
	await store.checkIsProxyRunning()
	await getCrashReport()
	loadingBar.finish()
	isProxyReloading.value = false
}

const isProxyStopping = ref(false)
async function stopProxy() {
	isProxyStopping.value = true
	loadingBar.start()
	await axios.get<string>(`${import.meta.env.VITE_APP_API}/api/stop`).then(() => {
		useToast().success('Proxy stopped')
	}).catch(err => {
		useToast().error(
			h(ToastDesc, {
				title: 'Proxy stop failed',
				message: err.response.data ?? err.message,
			}), {
			timeout: 5000,
		})
	})
	await store.checkIsProxyRunning()
	await getCrashReport()
	loadingBar.finish()
	isProxyStopping.value = false
}

const isCrashReportLoading = ref(false)
const crashReport = ref<HaproxyCrashReasonsData | null>(null)
async function getCrashReport() {
	isCrashReportLoading.value = true
	await axios.get<HaproxyCrashReasonsData>(`${import.meta.env.VITE_APP_API}/api/proxy/crash`).then((res) => {
		crashReport.value = res.data
	}).catch(err => {
		useToast().error(
			h(ToastDesc, {
				title: 'Crash report fetch failed',
				message: err.response.data ?? err.message,
			}), {
			timeout: 5000,
		})
	})
	isCrashReportLoading.value = false
}

</script>

<template>
	<n-card title="Home">
		<n-space vertical>
			<n-flex align="center">
				<n-button v-if="!store.isProxyRunning" :loading="isProxyStarting" @click="startProxy()" type="info">
					Start
				</n-button>
				<n-button v-if="store.isProxyRunning" :loading="isProxyStopping" @click="stopProxy()" type="error">
					Stop
				</n-button>
				<n-button @click="reloadProxy()" :loading="isProxyReloading" type="info">
					Reload
				</n-button>
			</n-flex>
			<n-h3>
				Status
			</n-h3>
			<n-space vertical>
				<n-flex align="center">
					<n-tag type="error" v-if="crashReport?.has_crashed">Crashed</n-tag>
					<n-tag type="success" v-if="!crashReport?.has_crashed">Healthy</n-tag>
				</n-flex>
				<n-alert v-if="crashReport?.address_in_use" title="Address/Port Already in use" type="warning">
					<n-code :show-line-numbers="true" :word-wrap="true"
						:code="crashReport?.address_in_use_log"></n-code>
					<n-space vertical>
						<div><strong>Option 2:</strong> Change the listening port</div>
					</n-space>
				</n-alert>
				<n-alert v-if="crashReport?.permission_denied_port" title="Permission denied using Port" type="warning">
					<n-code :show-line-numbers="true" :word-wrap="true"
						:code="crashReport?.permission_denied_port_log"></n-code>
					<n-space vertical>
						<div><strong>Option 1:</strong> Run the proxy with another user that has permission to bind to the port</div>
						<div><strong>Option 2:</strong> Change the listening port</div>
					</n-space>
				</n-alert>
			</n-space>
		</n-space>
	</n-card>
</template>