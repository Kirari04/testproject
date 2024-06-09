<script setup lang="ts">
import type { HaproxyLog } from 'env';
import { NCard, useLoadingBar, NInput, NButton, NSpace } from 'naive-ui'
import { onMounted, ref, h, watch, onUnmounted } from 'vue';
import { useToast } from 'vue-toastification';
import ToastDesc from '@/components/ToastDesc.vue'
import axios from 'axios';
const loadingBar = useLoadingBar()

onMounted(async () => {
	loadingBar.start()
	await getHaproxyLogs()
	setTimeout(() => {
		loadingBar.finish()
	}, 500)
})

let intv: any = null;
onMounted(() => {
	intv = setInterval(() => {
		getHaproxyLogs()
	}, 5 * 1000)
})
onUnmounted(() => {
	clearInterval(intv)
})

const logs = ref<HaproxyLog[]>([])
const rawLogs = ref<string>("")

watch(logs, (logs) => {
	rawLogs.value = logs.reverse().map(l => `${new Date(l.created_at).toLocaleString()} - ${l.data}`).join("");
	const txt = (document.querySelector("#haproxylogs textarea") as HTMLTextAreaElement)
	setTimeout(() => {
		txt.scrollTo({
			left: 0,
			top: txt.scrollHeight,
			behavior: 'smooth',
		})
	}, 300)
})

async function getHaproxyLogs() {
	await axios.get<HaproxyLog[]>(`${import.meta.env.VITE_APP_API}/api/haproxy/logs`).then(data => {
		logs.value = data.data
	}).catch(err => {
		useToast().error(
			h(ToastDesc, {
				title: 'Failed to get haproxy logs',
				message: err.response.data ?? err.message,
			}), {
			timeout: 5000,
		})
	})
}

async function reloadLogs() {
	loadingBar.start()
	await getHaproxyLogs()
	loadingBar.finish()
}
</script>

<template>
	<n-card title="Haproxy Logs">
		<n-space vertical>
			<n-button type="primary" @click="reloadLogs()">Refresh</n-button>
			<n-input id="haproxylogs" v-model:value="rawLogs" type="textarea" placeholder="Haproxy Logs"
				style="min-height: 200px; height: 600px; max-height: calc(100vh - 200px);" readonly />
		</n-space>
	</n-card>
</template>