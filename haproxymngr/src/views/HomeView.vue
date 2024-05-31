<script setup lang="ts">
import axios from 'axios';
import { NCard, useLoadingBar, NTag, NFlex, NButton } from 'naive-ui'
import { onMounted, ref } from 'vue';
import { useToast } from 'vue-toastification';
const loadingBar = useLoadingBar()
onMounted(async () => {
	loadingBar.start()
	await checkStatus()
	setTimeout(() => {
		loadingBar.finish()
	}, 500)
})

const isRunning = ref(false)
async function checkStatus() {
	await axios.get<string>(`${import.meta.env.VITE_APP_API}/api/status`)
		.then(res => {
			isRunning.value = res.data === 'ok'
		})
		.catch(err => {
			console.log(err.message)
		})
}

async function startProxy() {
	loadingBar.start()
	await axios.get<string>(`${import.meta.env.VITE_APP_API}/api/start`)
	useToast().success('Proxy started')
	await checkStatus()
	loadingBar.finish()
}

async function stopProxy() {
	loadingBar.start()
	await axios.get<string>(`${import.meta.env.VITE_APP_API}/api/stop`)
	useToast().success('Proxy stopped')
	await checkStatus()
	loadingBar.finish()
}

</script>

<template>
	<n-card title="Home">
		<n-flex align="center">
			<n-tag v-if="isRunning" type="success">
				Proxy is running
			</n-tag>
			<n-tag v-if="!isRunning" type="error">
				Proxy is off
			</n-tag>
			<n-button v-if="!isRunning" @click="startProxy()" type="info">
				Start
			</n-button>
			<n-button v-if="isRunning" @click="stopProxy()" type="error">
				Stop
			</n-button>
		</n-flex>
		<p>This is the home page</p>
	</n-card>
</template>