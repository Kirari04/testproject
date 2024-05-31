<script setup lang="ts">
import { h, nextTick, onMounted, ref, watch } from 'vue'
import { type Component } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'
import { NLayout, NLayoutSider, NIcon, NMenu, NLoadingBarProvider, NFlex } from 'naive-ui'
import type { MenuOption } from 'naive-ui'
import {
	HomeRound,
	InfoRound,
	WebRound,
} from '@vicons/material'
import { useStore } from './stores/store'
import axios from 'axios'


function renderIcon(icon: Component) {
	return () => h(NIcon, null, { default: () => h(icon) })
}

const store = useStore()
store.checkIsProxyRunning()
setInterval(() => {
	store.checkIsProxyRunning()
}, 10 * 1000)

const router = useRouter()
const collapsed = ref(false)
const activeKey = ref("-")

const menuOptions: MenuOption[] = [
	{
		label: 'Home',
		key: '/',
		icon: renderIcon(HomeRound),
	},
	{
		label: 'Http Proxy',
		key: '/http-proxy',
		icon: renderIcon(WebRound),
	},
	{
		label: 'About',
		key: '/about',
		icon: renderIcon(InfoRound),
	},
]

watch(activeKey, (key) => {
	if (key === "-") return
	router.push(key)
})

const animateRoute = ref(false)
router.beforeEach(async (to) => {
	animateRoute.value = true
	await new Promise((resolve) => setTimeout(resolve, 200))
})
router.afterEach((to) => {
	setTimeout(() => {
		animateRoute.value = false
	}, 200)
	activeKey.value = to.path
})
</script>

<template>
	<n-loading-bar-provider>
		<n-layout id="app-layout" has-sider>
			<n-layout-sider bordered collapse-mode="width" :collapsed-width="64" :width="240" :collapsed="collapsed"
				show-trigger @collapse="collapsed = true" @expand="collapsed = false">
				<n-menu v-model:value="activeKey" :collapsed="collapsed" :collapsed-width="64" :collapsed-icon-size="22"
					:options="menuOptions" />
			</n-layout-sider>
			<n-layout vertical>
				<div :class="`space ${animateRoute ? 'animate-in' : 'animate-out'}`">
					<n-flex size="large" vertical>
						<RouterView />
					</n-flex>
				</div>
			</n-layout>
		</n-layout>
	</n-loading-bar-provider>
</template>

<style>
#app-layout {
	height: 100vh;
}

.space {
	padding: 10px;
}

.animate-in {
	transition: all 0.2s ease;
	opacity: 0;
	transform: translateY(30px);
}

.animate-out {
	transition: all 0.2s ease;
	opacity: 1;
	transform: translateY(0px);
}
</style>