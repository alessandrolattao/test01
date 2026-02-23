<template>
  <audio v-if="src" :src="authSrc" controls class="h-8 w-48"></audio>
  <span v-else class="text-sm text-base-content/50">No audio</span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  src: { type: String, default: null },
})

const authSrc = computed(() => {
  if (!props.src) return null
  const token = localStorage.getItem('admin_token')
  return token ? `${props.src}?token=${token}` : props.src
})
</script>
