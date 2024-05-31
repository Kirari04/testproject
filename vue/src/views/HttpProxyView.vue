<script setup lang="ts">
import { NCard, useLoadingBar, NSpace, NInput, NInputNumber, NButton, NIcon, NTable, NTag } from 'naive-ui'
import { onMounted, ref, h } from 'vue'
import { type Component } from 'vue'
import {
	AddRound,
	MinusRound,
} from '@vicons/material'
import axios from 'axios';
import type { Frontend } from 'env';
import { useToast } from 'vue-toastification';
import { useStore } from '@/stores/store';
import ToastDesc from '@/components/ToastDesc.vue'
function renderIcon(icon: Component) {
	return () => h(NIcon, null, { default: () => h(icon) })
}

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

const port = ref(80)
const domain = ref('')
const backends = ref<{ addr: string }[]>([{ addr: '' }])

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
		backends: backends.value.map(b => ({ address: b.addr })),
	})
		.then(res => {
			useToast().success('Proxy created')
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
	port.value = 80
	domain.value = ''
	backends.value = [{ addr: '' }]
	await getProxies()
	loadingBar.finish()
}

async function deleteProxy(pr: Frontend) {
	loadingBar.start()
	await axios.delete<string>(`${import.meta.env.VITE_APP_API}/api/proxy`, {
		data: {
			id: pr.id
		}
	})
		.then(res => {
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
	<n-card title="Create Proxy">
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
				<n-space v-for="(backend, i) in backends">
					<n-input v-model:value="backend.addr" type="text" placeholder="127.0.0.1:8080" />
					<n-button :render-icon="renderIcon(MinusRound)" :disabled="backends.length === 1"
						@click="removeBackend(i)"></n-button>
				</n-space>
				<n-space>
					<n-button :render-icon="renderIcon(AddRound)" @click="addBackend()"></n-button>
				</n-space>
			</n-space>

			<n-button type="primary" @click="createProxy()">Save</n-button>
		</n-space>
	</n-card>
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
					<tr v-for="pr in proxies">
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