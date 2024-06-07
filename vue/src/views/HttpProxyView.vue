<script setup lang="ts">
import { NCard, useLoadingBar, NSpace, NInput, NInputNumber, NButton, NIcon, NTable, NTag, NSelect, NPopover } from 'naive-ui'
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
		.then(async res => {
			proxies.value = res.data
			getProxiesStatus()
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
async function getProxiesStatus() {
	await axios.get<FrontendStatus[]>(`${import.meta.env.VITE_APP_API}/api/proxies/status`)
		.then(res => {
			proxies_status.value = res.data
		})
		.catch(err => {
			useToast().error(
				h(ToastDesc, {
					title: 'Failed get proxies status',
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

function getProxyStatus(frontendId: number, backendId: number): ProxyStatus | null {
	const frontend = proxies_status.value
		.find(pr => pr.frontend_id === frontendId);
	if (!frontend) return null;
	const servers = frontend.Servers;
	const server = servers.find(s => s.server_id === backendId)
	if (!server) return null;
	for (const s of statuses) {
		if (server[s.code] === 1) {
			return s;
		}
	}
	return statuses[0];
}

function getProxyStats(frontendId: number): (FrontendStatus | null) {
	const frontend = proxies_status.value
		.find(pr => pr.frontend_id === frontendId);
	if (!frontend) return null;
	return frontend;
}

type ProxyCodes = "hana" | "sockerr" | "l4ok" | "l4tout" | "l4con" | "l6ok" | "l6tout" | "l6rsp" | "l7tout" | "l7rsp" | "l7ok" | "l7okc" | "l7sts" | "procerr" | "proctout" | "procok";
type ProxyStatus = { code: ProxyCodes; description: string; isError: boolean; };
const statuses: ProxyStatus[] = [
		{ code: "sockerr", description: "socket error", isError: true },
		{ code: "l4ok", description: "check passed on layer 4, no upper layers testing enabled", isError: false },
		{ code: "l4tout", description: "layer 1-4 timeout" , isError: true },
		{ code: "l4con", description: "layer 1-4 connection problem, for example 'Connection refused' (tcp rst) or 'No route to host' (icmp)", isError: true },
		{ code: "l6ok", description: "check passed on layer 6", isError: false },
		{ code: "l6tout", description: "layer 6 (SSL) timeout", isError: true },
		{ code: "l6rsp", description: "layer 6 invalid response - protocol error", isError: true },
		{ code: "l7tout", description: "layer 7 (HTTP/SMTP) timeout", isError: true },
		{ code: "l7rsp", description: "layer 7 invalid response - protocol error", isError: true },
		{ code: "l7ok", description: "check passed on layer 7" , isError: false },
		{ code: "l7okc", description: "check conditionally passed on layer 7, for example 404 with disable-on-404", isError: false },
		{ code: "l7sts", description: "layer 7 response error, for example HTTP 5xx", isError: true },
	];


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
						<th>Requests</th>
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
									<a :href="`http://${pr.domain}:${pr.port}`" target="_blank"
										rel="noopener noreferrer">
										{{ pr.domain }}
									</a>
								</li>
								<li v-for="alias in pr.aliases">
									<a :href="`http://${alias.domain}:${pr.port}`" target="_blank"
										rel="noopener noreferrer">
										{{ alias.domain }}
									</a>
								</li>
							</ul>

						</td>
						<td>
							<n-space vertical>
								<NTag type="info">HTTP 1xx - {{ getProxyStats(pr.id)?.responses_total_1xx }}</NTag>
								<NTag type="success">HTTP 2xx - {{ getProxyStats(pr.id)?.responses_total_2xx }}</NTag>
								<NTag type="success">HTTP 3xx - {{ getProxyStats(pr.id)?.responses_total_3xx }}</NTag>
								<NTag type="warning">HTTP 4xx - {{ getProxyStats(pr.id)?.responses_total_4xx }}</NTag>
								<NTag type="error">HTTP 5xx - {{ getProxyStats(pr.id)?.responses_total_5xx }}</NTag>
								<NTag>Other - {{ getProxyStats(pr.id)?.responses_total_other }}</NTag>
							</n-space>
						</td>
						<td>
							<NTable size="small">
				<tbody>
					<tr>
						<th>Address</th>
						<th>Status</th>
					</tr>
					<tr v-for="b in pr.backends" :key="b.id">
						<td>
							{{ b.address }}
						</td>
						<td>
							<n-popover trigger="hover">
								<template #trigger>
									<n-tag :type="getProxyStatus(pr.id, b.id)?.isError ? 'error' : 'success'">
										{{ getProxyStatus(pr.id, b.id)?.code }}
									</n-tag>
								</template>
								<span>{{ getProxyStatus(pr.id, b.id)?.description }}</span>
							</n-popover>
						</td>
					</tr>
				</tbody>
				</NTable>
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