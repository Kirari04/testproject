<script setup lang="ts">
import { NCard, useLoadingBar, NSpace, NButton, NTable, NTag, NPopover } from 'naive-ui'
import { onMounted, ref, h, onUnmounted } from 'vue'
import axios from 'axios';
import type { Frontend, FrontendStatus } from 'env';
import { useToast } from 'vue-toastification';
import { useStore } from '@/stores/store';
import ToastDesc from '@/components/ToastDesc.vue'
import CreateProxy from '@/components/CreateProxy.vue'

let intv: any = null;

const loadingBar = useLoadingBar()
onMounted(async () => {
	loadingBar.start()
	await getProxies()
	intv = setInterval(() => {
		getProxiesStatus()
	}, 4000)
	setTimeout(() => {
		loadingBar.finish()
	}, 500)
})

onUnmounted(() => {
	clearInterval(intv)
})

const store = useStore()

const proxies = ref<Frontend[]>([])
const proxies_status = ref<FrontendStatus[]>([])

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
					message: err.response.data ?? err.message,
				}), {
				timeout: 5000,
			})
		})

	await getProxies()
	loadingBar.finish()
}

async function reloadProxies() {
	loadingBar.start()
	await getProxies()
	loadingBar.finish()
}

async function getProxies() {
	await axios.get<Frontend[]>(`${import.meta.env.VITE_APP_API}/api/proxies`)
		.then(async res => {
			proxies.value = res.data
			getProxiesStatus()
		})
		.catch(err => {
			useToast().error(
				h(ToastDesc, {
					title: 'Failed get proxies',
					message: err.response.data ?? err.message,
				}), {
				timeout: 5000,
			})
		})
}
async function getProxiesStatus() {
	if (!store.isProxyRunning) {
		return
	}
	await axios.get<FrontendStatus[]>(`${import.meta.env.VITE_APP_API}/api/proxies/status`)
		.then(res => {
			proxies_status.value = res.data
		})
		.catch(err => {
			useToast().error(
				h(ToastDesc, {
					title: 'Failed get proxies status',
					message: err.response.data ?? err.message,
				}), {
				timeout: 5000,
			})
		})
}

function getProxyStatus(frontendId: number, backendId: number): ProxyStatus | null {
	const frontend = proxies_status.value
		.find(pr => pr.frontend_id === frontendId);
	if (!frontend) return statuses[0];
	const servers = frontend.Servers;
	const server = servers.find(s => s.server_id === backendId)
	if (!server) return statuses[0];
	for (const s of statuses) {
		if ((server as any)[s.code as any] === 1) {
			return s;
		}
	}
	return statuses[0];
}

type ProxyCodes = "unknown" | "sockerr" | "l4ok" | "l4tout" | "l4con" | "l6ok" | "l6tout" | "l6rsp" | "l7tout" | "l7rsp" | "l7ok" | "l7okc" | "l7sts";
type ProxyStatus = { code: ProxyCodes; description: string; isError: boolean; };
const statuses: ProxyStatus[] = [
	{ code: "unknown", description: "this error happens if now stats are available", isError: true },
	{ code: "sockerr", description: "socket error", isError: true },
	{ code: "l4ok", description: "check passed on layer 4, no upper layers testing enabled", isError: false },
	{ code: "l4tout", description: "layer 1-4 timeout", isError: true },
	{ code: "l4con", description: "layer 1-4 connection problem, for example 'Connection refused' (tcp rst) or 'No route to host' (icmp)", isError: true },
	{ code: "l6ok", description: "check passed on layer 6", isError: false },
	{ code: "l6tout", description: "layer 6 (SSL) timeout", isError: true },
	{ code: "l6rsp", description: "layer 6 invalid response - protocol error", isError: true },
	{ code: "l7tout", description: "layer 7 (HTTP/SMTP) timeout", isError: true },
	{ code: "l7rsp", description: "layer 7 invalid response - protocol error", isError: true },
	{ code: "l7ok", description: "check passed on layer 7", isError: false },
	{ code: "l7okc", description: "check conditionally passed on layer 7, for example 404 with disable-on-404", isError: false },
	{ code: "l7sts", description: "layer 7 response error, for example HTTP 5xx", isError: true },
];


</script>

<template>
	<n-card title="Proxies">
		<n-space vertical>
			<n-space>
				<n-button @click="reloadProxies()">Reload List</n-button>
				<CreateProxy @onCreated="reloadProxies" />
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
							<ul>
								<li>
									<a :href="`${pr.https ? 'https' : 'http'}://${pr.domain}:${pr.port}`"
										target="_blank" rel="noopener noreferrer">
										{{ pr.domain }}
									</a>
								</li>
								<li v-for="alias in pr.aliases">
									<a :href="`${pr.https ? 'https' : 'http'}://${alias.domain}:${pr.port}`"
										target="_blank" rel="noopener noreferrer">
										{{ alias.domain }}
									</a>
								</li>
							</ul>

						</td>
						<td>
							<ul>
								<li v-for="b in pr.backends" :key="b.id">

									<n-popover trigger="hover">
										<template #trigger>
											<n-tag :type="getProxyStatus(pr.id, b.id)?.isError ? 'error' : 'success'">
												{{ getProxyStatus(pr.id, b.id)?.code }}
											</n-tag>
										</template>
										<span>{{ getProxyStatus(pr.id, b.id)?.description }}</span>
									</n-popover>
									<n-popover trigger="hover">
										<template #trigger>
											<n-tag>
												{{ b.https ? 'https' : 'http' }}
											</n-tag>
										</template>
										<span>
											The backend is being requested over
											<a :href="`${b.https ? 'https' : 'http'}://${b.address}`" target="_blank"
												rel="noopener noreferrer">
												{{ b.https ? 'https' : 'http' }}://{{ b.address }}
											</a>
										</span>
									</n-popover>
									<n-tag>
										{{ b.address }}
									</n-tag>
								</li>
							</ul>
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