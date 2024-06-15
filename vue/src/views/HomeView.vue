<script setup lang="ts">
import axios from 'axios';
import { NCard, useLoadingBar, NTag, NFlex, NButton, NSpace } from 'naive-ui'
import { onMounted, ref, h } from 'vue'
import { useToast } from 'vue-toastification';
import { useStore } from '@/stores/store'
import ToastDesc from '@/components/ToastDesc.vue'

const loadingBar = useLoadingBar()
onMounted(async () => {
	loadingBar.start()
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
	loadingBar.finish()
	isProxyStopping.value = false
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
		</n-space>
	</n-card>
</template>