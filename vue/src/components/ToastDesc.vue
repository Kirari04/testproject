<template>
    <div class="container">
        <div><strong>{{ title }}</strong></div>
        <div>{{ message }}</div>
    </div>
</template>

<script lang="ts" setup>

defineProps<{
    title: string
    message: string
}>()
</script>