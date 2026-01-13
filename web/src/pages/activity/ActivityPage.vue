<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()
const gameKey = String(route.params.gameKey || '')
const id = String(route.params.id || '')
</script>

<template>
  <div class="mx-auto w-full max-w-6xl px-4 py-4">
    <div class="text-lg font-semibold">活动详情</div>
    <div class="mt-2 text-sm text-muted-foreground">gameKey={{ gameKey }} · id={{ id }}</div>
  </div>
</template>

