<template>
  <AdminLayout>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">Candidates</h1>
      <div v-if="data" class="badge badge-lg">{{ data.length }} total</div>
    </div>

    <!-- Stats -->
    <div v-if="data && data.length > 0" class="stats shadow mb-6 w-full">
      <div class="stat">
        <div class="stat-title">Total Candidates</div>
        <div class="stat-value">{{ data.length }}</div>
      </div>
      <div class="stat">
        <div class="stat-title">Average Score</div>
        <div class="stat-value">{{ averageScore }}</div>
      </div>
      <div class="stat">
        <div class="stat-title">Completed</div>
        <div class="stat-value">{{ completedCount }}</div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="isPending" class="flex justify-center p-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="alert alert-error">
      <span>Failed to load candidates.</span>
      <button class="btn btn-sm btn-ghost" @click="refresh()">Retry</button>
    </div>

    <!-- Empty -->
    <div v-else-if="data && data.length === 0" class="text-center py-12 text-base-content/50">
      <p class="text-lg">No candidates yet</p>
    </div>

    <!-- Table -->
    <div v-else-if="data" class="overflow-x-auto">
      <table class="table table-zebra">
        <thead>
          <tr>
            <th>#</th>
            <th>Name</th>
            <th>Email</th>
            <th>Score</th>
            <th>AI Score</th>
            <th>Audio</th>
            <th>Date</th>
            <th>Actions</th>
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
              <span v-else-if="candidate.analysis_status === 'failed'" class="badge badge-sm badge-error">Failed</span>
              <span v-else-if="!candidate.completed" class="badge badge-sm badge-ghost">Incomplete</span>
              <span v-else class="text-sm text-base-content/30">-</span>
            </td>
            <td>
              <AudioPlayer
                v-if="candidate.audio_path"
                :src="candidateAudioUrl(candidate.id)"
              />
              <span v-else class="text-sm text-base-content/50">No audio</span>
            </td>
            <td class="text-sm">{{ formatDate(candidate.created_at) }}</td>
            <td>
              <RouterLink
                :to="{ name: 'admin-candidate-detail', params: { id: candidate.id } }"
                class="btn btn-ghost btn-sm"
              >
                View
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

function aiScoreBadge(score) {
  if (score >= 70) return 'badge-success'
  if (score >= 40) return 'badge-warning'
  return 'badge-error'
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  }) + ' ' + d.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>
