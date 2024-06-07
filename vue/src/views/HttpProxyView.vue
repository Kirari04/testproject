<script setup lang="ts">
import { NCard, useLoadingBar, NSpace, NInput, NInputNumber, NButton, NIcon, NTable, NTag, NSelect } from 'naive-ui'
import { onMounted, ref, h } from 'vue'
import axios from 'axios';
import type { Frontend } from 'env';
import { useToast } from 'vue-toastification';
import { useStore } from '@/stores/store';
import ToastDesc from '@/components/ToastDesc.vue'
import CreateProxy from '@/components/CreateProxy.vue'

const loadingBar = useLoadingBar()
onMounted(async () => {
	loadingBar.start()
	await getProxies()
	setTimeout(() => {
		loadingBar.finish()
	}, 500)
})

const store = useStore()

const proxies = ref<Frontend[]>([])

async function deleteProxy(pr: Frontend) {
	loadingBar.start()
	await axios.delete<string>(`${import.meta.env.VITE_APP_API}/api/proxy`, {
		data: {
			id: pr.id
		}
	})
		.then(() => {
			useToast().success('Proxy deleted')
		})
		.catch(err => {
			useToast().error(
				h(ToastDesc, {
					title: 'Failed to delete proxy',
					message: err.message,
				}), {
				timeout: 5000,
			})
		})

	await getProxies()
	loadingBar.finish()
}

async function getProxies() {
	await axios.get<Frontend[]>(`${import.meta.env.VITE_APP_API}/api/proxies`)
		.then(res => {
			proxies.value = res.data
		})
		.catch(err => {
			useToast().error(
				h(ToastDesc, {
					title: 'Failed get proxies',
					message: err.message,
				}), {
				timeout: 5000,
			})
		})
}

async function runApply() {
	loadingBar.start()
	await axios.post<string>(`${import.meta.env.VITE_APP_API}/api/proxy/apply`)
		.then(res => {
			console.log(res.data)
			useToast().success('Proxy applied')
		})
		.catch(err => {
			useToast().error(
				h(ToastDesc, {
					title: 'Failed to apply proxy',
					message: err.message,
				}), {
				timeout: 5000,
			})
		})
	store.checkIsProxyRunning()
	loadingBar.finish()
}
</script>

<template>
	<n-card title="Proxies">
		<n-space vertical>
			<n-space>
				<n-tag v-if="store.isProxyRunning" type="success">
					Proxy is running
				</n-tag>
				<n-tag v-if="!store.isProxyRunning" type="error">
					Proxy is off
				</n-tag>
				<n-button type="primary" @click="runApply()">Apply</n-button>
				<CreateProxy />
			</n-space>
			<n-table :single-line="false">
				<thead>
					<tr>
						<th>ID</th>
						<th>Listen</th>
						<th>Domain</th>
						<th>Backends</th>
						<th>Action</th>
					</tr>
				</thead>
				<tbody>
					<tr v-for="pr in proxies" :key="pr.id">
						<td>
							{{ pr.id }}
						</td>
						<td>
							:{{ pr.port }}
						</td>
						<td>
							<a :href="`http://${pr.domain}:${pr.port}`" target="_blank" rel="noopener noreferrer">
								{{ pr.domain }}
							</a>
							<span v-for="alias in pr.aliases">
								<br>
								<a :href="`http://${alias.domain}:${pr.port}`" target="_blank"
									rel="noopener noreferrer">
									{{ alias.domain }}
								</a>
							</span>
						</td>
						<td>
							{{ pr.backends.map(b => b.address).join(', ') }}
						</td>
						<td>
							<n-space>
								<n-button type="error" @click="deleteProxy(pr)">Delete</n-button>
							</n-space>
						</td>
					</tr>
				</tbody>
			</n-table>
		</n-space>
	</n-card>
</template>