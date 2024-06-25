<script setup lang="ts">
import axios from 'axios';
import { NCard, useLoadingBar, NTag, NFlex, NButton, NSpace, NInputNumber, NInput } from 'naive-ui'
import { onMounted, ref, h } from 'vue'
import { useToast } from 'vue-toastification';
import { useStore } from '@/stores/store'
import ToastDesc from '@/components/ToastDesc.vue'
import type { Settings } from 'env';

const loadingBar = useLoadingBar()
onMounted(async () => {
	loadingBar.start()
	setTimeout(() => {
		loadingBar.finish()
	}, 500)
})

const store = useStore()

onMounted(async () => {
	loadSettings()
})

const isLoading = ref(false)
const acmeEmail = ref("")
const acmeCloudflareDNSAPIToken = ref("")

async function saveSettings() {
	isLoading.value = true
	await axios.post<string>(`${import.meta.env.VITE_APP_API}/api/settings`, {
		acme_email: acmeEmail.value,
		acme_cloudflare_dns_api_token: acmeCloudflareDNSAPIToken.value,
	})
		.then(() => {
			useToast().success('Settings saved')
		})
		.catch(err => {
			useToast().error(
				h(ToastDesc, {
					title: 'Settings save failed',
					message: err.response.data ?? err.message,
				}), {
				timeout: 5000,
			})
		})
	isLoading.value = false
	loadSettings()
}

async function loadSettings() {
	loadingBar.start()
	isLoading.value = true
	await axios.get<Settings>(`${import.meta.env.VITE_APP_API}/api/settings`)
		.then(res => {
			acmeEmail.value = res.data.acme_email
			acmeCloudflareDNSAPIToken.value = res.data.acme_cloudflare_dns_api_token
		})
		.catch(err => {
			useToast().error(
				h(ToastDesc, {
					title: 'Settings load failed',
					message: err.response.data ?? err.message,
				}), {
				timeout: 5000,
			})
		})
	loadingBar.finish()
	isLoading.value = false
}

</script>

<template>
	<n-card title="Settings">
		<n-space vertical>
			<n-flex align="center">
				<!-- <n-button v-if="!store.isProxyRunning" :loading="isProxyStarting" @click="startProxy()" type="info">
					Start
				</n-button> -->
			</n-flex>
		</n-space>
		<n-card>
			<h3>Acme (Certificate Requests)</h3>
			<form autocomplete="off" data-lpignore="true">
				<n-space size="large" vertical>
					<n-space vertical>
						E-Mail
						<n-space>
							<n-input placeholder="example@example.com" v-model:value="acmeEmail" />
						</n-space>
					</n-space>
					<n-space vertical>
						Cloudflare DNS API Token
						<n-space>
							<n-input :showPasswordToggle="true" placeholder="xxx" type="password"
								v-model:value="acmeCloudflareDNSAPIToken" />
						</n-space>
					</n-space>
					<n-button :loading="isLoading" type="primary" @click="saveSettings()">Save</n-button>
				</n-space>
			</form>
			<!-- <h3>Hard limit</h3>
			<n-space>
				<div>
					Limit
					<n-input-number placeholder="0" />
				</div>
				<div>
					Period (seconds)
					<n-input-number placeholder="0" />
				</div>
			</n-space> -->
		</n-card>
	</n-card>
</template>