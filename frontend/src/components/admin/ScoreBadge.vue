<template>
  <div v-if="score === null || score === undefined" class="badge badge-ghost">Pending</div>
  <div v-else class="badge" :class="badgeClass">{{ score }}/{{ maxScore }}</div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  score: { type: Number, default: null },
  maxScore: { type: Number, default: 100 },
})

const badgeClass = computed(() => {
  if (props.score === null || props.score === undefined) return 'badge-ghost'
  const pct = (props.score / props.maxScore) * 100
  if (pct >= 80) return 'badge-success'
  if (pct >= 50) return 'badge-warning'
  return 'badge-error'
})
</script>
