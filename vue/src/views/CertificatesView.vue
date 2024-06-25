<script setup lang="ts">
import axios from 'axios';
import { NCard, useLoadingBar, NTag, NFlex, NButton, NSpace, NTable, NPopover } from 'naive-ui'
import { onMounted, h, ref } from 'vue'
import ToastDesc from '@/components/ToastDesc.vue'
import { useToast } from 'vue-toastification';
import { useStore } from '@/stores/store'
import UploadCertificate from '@/components/UploadCertificate.vue';
import RequestCertificate from '@/components/RequestCertificate.vue';
import type { Certificate } from 'env';

const loadingBar = useLoadingBar()
onMounted(async () => {
	loadingBar.start()
	await getCertificates()
	setTimeout(() => {
		loadingBar.finish()
	}, 500)
})

const store = useStore()

const cerificates = ref<Certificate[]>([])

async function getCertificates() {
	await axios.get<Certificate[]>(`${import.meta.env.VITE_APP_API}/api/certificates`)
		.then(async res => {
			cerificates.value = res.data
		})
		.catch(err => {
			useToast().error(
				h(ToastDesc, {
					title: 'Failed get certificates',
					message: err.response.data ?? err.message,
				}), {
				timeout: 5000,
			})
		})
}

async function reloadCertificates() {
	loadingBar.start()
	await getCertificates()
	loadingBar.finish()
}

async function deleteCertificate(cert: Certificate) {
	loadingBar.start()
	await axios.delete<string>(`${import.meta.env.VITE_APP_API}/api/certificate`, {
		data: {
			id: cert.id
		}
	})
		.then(() => {
			useToast().success('Certificate deleted')
		})
		.catch(err => {
			useToast().error(
				h(ToastDesc, {
					title: 'Failed to delete certificate',
					message: err.response.data ?? err.message,
				}), {
				timeout: 5000,
			})
		})

	await getCertificates()
	loadingBar.finish()
}
</script>

<template>
	<n-card title="Certificates">
		<n-space vertical>
			<n-space>
				<n-button @click="reloadCertificates()">Reload List</n-button>
				<UploadCertificate @on-created="reloadCertificates()" />
				<RequestCertificate @on-created="reloadCertificates()" />
			</n-space>
			<n-table :single-line="false">
				<thead>
					<tr>
						<th>ID</th>
						<th>Created</th>
						<th>Updated</th>
						<th>Name</th>
						<th>File</th>
						<th>Action</th>
					</tr>
				</thead>
				<tbody>
					<tr v-for="cert in cerificates" :key="cert.id">
						<td>
							{{ cert.id }}
						</td>
						<td>
							{{ new Date(cert.created_at).toLocaleString() }}
						</td>
						<td>
							{{ new Date(cert.updated_at).toLocaleString() }}
						</td>
						<td>
							{{ cert.name }}
						</td>
						<td>
							{{ cert.pem_path }}
						</td>
						<td>
							<n-button type="error" @click="deleteCertificate(cert)">Delete</n-button>
						</td>
					</tr>
				</tbody>
			</n-table>
		</n-space>
	</n-card>
</template>