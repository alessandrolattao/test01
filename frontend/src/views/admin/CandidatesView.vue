<template>
  <AdminLayout>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">{{ $t('admin.candidates.title') }}</h1>
      <div v-if="data" class="badge badge-lg">{{ $t('admin.candidates.total', { count: data.length }) }}</div>
    </div>

    <!-- Stats -->
    <div v-if="data && data.length > 0" class="stats shadow mb-6 w-full">
      <div class="stat">
        <div class="stat-title">{{ $t('admin.candidates.statTotal') }}</div>
        <div class="stat-value">{{ data.length }}</div>
      </div>
      <div class="stat">
        <div class="stat-title">{{ $t('admin.candidates.statAvg') }}</div>
        <div class="stat-value">{{ averageScore }}</div>
      </div>
      <div class="stat">
        <div class="stat-title">{{ $t('admin.candidates.statCompleted') }}</div>
        <div class="stat-value">{{ completedCount }}</div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="isPending" class="flex justify-center p-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="alert alert-error">
      <span>{{ $t('admin.candidates.loadError') }}</span>
      <button class="btn btn-sm btn-ghost" @click="refresh()">{{ $t('common.retry') }}</button>
    </div>

    <!-- Empty -->
    <div v-else-if="data && data.length === 0" class="text-center py-12 text-base-content/50">
      <p class="text-lg">{{ $t('admin.candidates.empty') }}</p>
    </div>

    <!-- Table -->
    <div v-else-if="data" class="overflow-x-auto">
      <table class="table table-zebra">
        <thead>
          <tr>
            <th>#</th>
            <th>{{ $t('admin.candidates.colName') }}</th>
            <th>{{ $t('admin.candidates.colEmail') }}</th>
            <th>{{ $t('admin.candidates.colScore') }}</th>
            <th>{{ $t('admin.candidates.colAiScore') }}</th>
            <th>{{ $t('admin.candidates.colAudio') }}</th>
            <th>{{ $t('admin.candidates.colDate') }}</th>
            <th>{{ $t('admin.candidates.colActions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(candidate, i) in data" :key="candidate.id" class="hover">
            <td>{{ i + 1 }}</td>
            <td class="font-medium">{{ candidate.first_name }} {{ candidate.last_name }}</td>
            <td class="text-sm">{{ candidate.email }}</td>
            <td>
              <ScoreBadge :score="candidate.total_score > 0 ? candidate.total_score : null" :max-score="maxScore" />
            </td>
            <td>
              <ScoreBadge v-if="candidate.ai_score != null" :score="candidate.ai_score" :max-score="100" />
              <span v-else-if="candidate.analysis_status === 'transcribing' || candidate.analysis_status === 'analyzing'" class="loading loading-dots loading-xs"></span>
              <span v-else-if="candidate.analysis_status === 'failed'" class="badge badge-sm badge-error">{{ $t('admin.candidates.failed') }}</span>
              <span v-else-if="!candidate.completed" class="badge badge-sm badge-ghost">{{ $t('admin.candidates.incomplete') }}</span>
              <span v-else class="text-sm text-base-content/30">-</span>
            </td>
            <td>
              <AudioPlayer
                v-if="candidate.audio_path"
                :src="candidateAudioUrl(candidate.id)"
              />
              <span v-else class="text-sm text-base-content/50">{{ $t('admin.candidates.noAudio') }}</span>
            </td>
            <td class="text-sm">{{ formatDate(candidate.created_at) }}</td>
            <td>
              <RouterLink
                :to="{ name: 'admin-candidate-detail', params: { id: candidate.id } }"
                class="btn btn-ghost btn-sm"
              >
                {{ $t('admin.candidates.view') }}
              </RouterLink>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </AdminLayout>
</template>

<script setup>
import { computed, onMounted, onUnmounted } from 'vue'
import { RouterLink } from 'vue-router'
import AdminLayout from '../../components/admin/AdminLayout.vue'
import AudioPlayer from '../../components/admin/AudioPlayer.vue'
import ScoreBadge from '../../components/admin/ScoreBadge.vue'
import { useCandidatesList, candidateAudioUrl } from '../../composables/useAdmin.js'

const { data, isPending, error, refresh } = useCandidatesList()

// Polling: refresh every 10s to catch new candidates and analysis updates
let pollInterval = null

onMounted(() => {
  pollInterval = setInterval(() => refresh(), 10000)
})

onUnmounted(() => {
  if (pollInterval) {
    clearInterval(pollInterval)
    pollInterval = null
  }
})

const maxScore = computed(() => {
  if (!data.value || data.value.length === 0) return 100
  const scores = data.value.filter((c) => c.completed).map((c) => c.total_score)
  return scores.length > 0 ? Math.max(...scores, 100) : 100
})

const averageScore = computed(() => {
  if (!data.value) return 0
  const completed = data.value.filter((c) => c.completed)
  if (completed.length === 0) return 0
  const total = completed.reduce((sum, c) => sum + c.total_score, 0)
  return Math.round(total / completed.length)
})

const completedCount = computed(() => {
  if (!data.value) return 0
  return data.value.filter((c) => c.completed).length
})

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  }) + ' ' + d.toLocaleTimeString(undefined, {
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>
