<script setup lang="ts">
import { NCard, useLoadingBar } from 'naive-ui'
import { onMounted } from 'vue';
const loadingBar = useLoadingBar()
onMounted(() => {
	loadingBar.start()
	setTimeout(() => {
		loadingBar.finish()
	}, 500)
})
</script>

<template>
	<n-card title="About">
		<h3>Welcome to HaproxyMngr</h3>

		<p>
			HaproxyMngr is an open-source project that aims to make configuring and managing a single instance of
			HAProxy
			easier than ever. We believe that with the right tools, anyone can unlock the full potential of HAProxy and
			bring their web applications to the next level.
		</p>

		<h3>What is HAProxy?</h3>

		<p>
			HAProxy is a powerful open-source load balancer and reverse proxy that has been widely adopted by
			organizations
			of all sizes. With its ability to handle large volumes of traffic and provide advanced features like session
			persistence, caching, and security, it's no wonder why HAProxy has become the go-to solution for many web
			applications.
		</p>

		<h3>The Problem</h3>

		<p>
			While HAProxy is an incredibly powerful tool, configuring it can be a daunting task, especially for those
			without extensive experience in networking or Linux systems administration. That's where we come in.
		</p>

		<h3>Our Solution</h3>

		<p>
			We created HaproxyMngr to provide a user-friendly graphical interface for managing a single instance of
			HAProxy.
			With our GUI, you can easily configure and manage your HAProxy instance, define frontends and backends, set
			up
			load balancing and routing rules, and monitor performance metrics - all without having to write a single
			line of
			code.
		</p>

		<h3>Important Note</h3>

		<p>
			HaproxyMngr is not affiliated with the HAProxy project or its maintainers. We simply created this tool to
			make
			it easier for users to manage their own HAProxy instances.
		</p>

		<h3>
			Why Open-Source?
		</h3>

		<p>
			At HaproxyMngr, we believe that open-source is the best way to bring people together around a shared goal.
			By
			releasing our project under an open-source license, we hope to create a community of like-minded individuals
			who
			can contribute to and benefit from our work.
		</p>
	</n-card>
</template>