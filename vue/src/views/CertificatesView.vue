<script setup lang="ts">
import axios from 'axios';
import { NCard, useLoadingBar, NTag, NFlex, NButton,  NSpace } from 'naive-ui'
import { onMounted } from 'vue'
import { useToast } from 'vue-toastification';
import { useStore } from '@/stores/store'
import UploadCertificate from '@/components/UploadCertificate.vue';

const loadingBar = useLoadingBar()
onMounted(async () => {
	loadingBar.start()
	setTimeout(() => {
		loadingBar.finish()
	}, 500)
})

const store = useStore()


</script>

<template>
	<n-card title="Certificates">
		<n-space vertical>
			<n-flex align="center">
				<UploadCertificate />
			</n-flex>
		</n-space>
	</n-card>
</template>