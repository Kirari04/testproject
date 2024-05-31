<script setup lang="ts">
import axios from 'axios';
import { NCard, useLoadingBar, NTag, NFlex, NButton } from 'naive-ui'
import { onMounted, ref } from 'vue';
import { useToast } from 'vue-toastification';
import { useStore } from '../stores/store'
const loadingBar = useLoadingBar()
onMounted(async () => {
	loadingBar.start()
	setTimeout(() => {
		loadingBar.finish()
	}, 500)
})

const store = useStore()

async function startProxy() {
	loadingBar.start()
	await axios.get<string>(`${import.meta.env.VITE_APP_API}/api/start`)
	useToast().success('Proxy started')
	await store.checkIsProxyRunning()
	loadingBar.finish()
}

async function stopProxy() {
	loadingBar.start()
	await axios.get<string>(`${import.meta.env.VITE_APP_API}/api/stop`)
	useToast().success('Proxy stopped')
	await store.checkIsProxyRunning()
	loadingBar.finish()
}

</script>

<template>
	<n-card title="Home">
		<n-flex align="center">
			<n-tag v-if="store.isProxyRunning" type="success">
				Proxy is running
			</n-tag>
			<n-tag v-if="!store.isProxyRunning" type="error">
				Proxy is off
			</n-tag>
			<n-button v-if="!store.isProxyRunning" @click="startProxy()" type="info">
				Start
			</n-button>
			<n-button v-if="store.isProxyRunning" @click="stopProxy()" type="error">
				Stop
			</n-button>
		</n-flex>
		<p>This is the home page</p>
	</n-card>
</template>