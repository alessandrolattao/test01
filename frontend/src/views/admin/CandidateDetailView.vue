<template>
  <AdminLayout>
    <!-- Loading -->
    <div v-if="isPending" class="flex justify-center items-center min-h-[50vh]">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="alert alert-error">
      <span>Failed to load candidate details.</span>
      <button class="btn btn-sm btn-ghost" @click="refresh()">Retry</button>
    </div>

    <template v-else-if="data">
      <!-- Breadcrumb -->
      <div class="breadcrumbs text-sm mb-6">
        <ul>
          <li><RouterLink to="/admin/candidates">Candidates</RouterLink></li>
          <li>{{ data.first_name }} {{ data.last_name }}</li>
        </ul>
      </div>

      <!-- Candidate info -->
      <div class="card bg-base-100 shadow-md mb-6">
        <div class="card-body">
          <div class="flex items-center justify-between flex-wrap gap-4">
            <div>
              <h2 class="card-title text-2xl">{{ data.first_name }} {{ data.last_name }}</h2>
              <p class="text-base-content/70">{{ data.email }}</p>
              <p class="text-sm text-base-content/50 mt-1">Applied on {{ formatDate(data.created_at) }}</p>
            </div>
            <div class="text-right">
              <div class="text-sm text-base-content/50">Total Score</div>
              <div class="text-4xl font-bold" :class="scoreColor">{{ data.total_score }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Audio -->
      <div class="card bg-base-100 shadow-md mb-6">
        <div class="card-body">
          <h3 class="card-title text-lg mb-3">Audio Presentation</h3>
          <audio v-if="data.audio_path" :src="audioSrc" controls class="w-full"></audio>
          <p v-else class="text-base-content/50">No audio recorded</p>
        </div>
      </div>

      <!-- Transcript -->
      <div v-if="data.audio_path" class="card bg-base-100 shadow-md mb-6">
        <div class="card-body">
          <h3 class="card-title text-lg mb-3">Transcript</h3>

          <!-- Loading state -->
          <div v-if="isTranscribing" class="flex items-center gap-3">
            <span class="loading loading-dots loading-md"></span>
            <span class="text-base-content/70">Transcription in progress...</span>
          </div>

          <!-- Transcript text -->
          <p v-else-if="data.transcript" class="whitespace-pre-wrap leading-relaxed">{{ data.transcript }}</p>

          <!-- Failed -->
          <div v-else-if="data.analysis_status === 'failed'" class="flex items-center justify-between">
            <span class="text-error">Transcription failed</span>
            <button class="btn btn-sm btn-outline btn-error" :disabled="reanalyzeIsPending" @click="handleReanalyze">
              <span v-if="reanalyzeIsPending" class="loading loading-spinner loading-xs"></span>
              Retry
            </button>
          </div>

          <!-- Pending -->
          <div v-else class="flex items-center gap-3">
            <span class="loading loading-dots loading-sm"></span>
            <span class="text-base-content/50">Waiting to start...</span>
          </div>
        </div>
      </div>

      <!-- AI Analysis -->
      <div v-if="data.audio_path" class="card bg-base-100 shadow-md mb-6">
        <div class="card-body">
          <div class="flex items-center justify-between mb-3">
            <h3 class="card-title text-lg">AI Analysis</h3>
            <div v-if="parsedAnalysis" class="flex items-center gap-2">
              <span class="text-sm text-base-content/50">AI Score</span>
              <div class="badge badge-lg font-bold" :class="aiScoreBadgeClass">
                {{ data.ai_score }}/100
              </div>
            </div>
          </div>

          <!-- Loading state -->
          <div v-if="isAnalyzing" class="flex items-center gap-3">
            <span class="loading loading-dots loading-md"></span>
            <span class="text-base-content/70">AI analysis in progress...</span>
          </div>

          <!-- Analysis result -->
          <div v-else-if="parsedAnalysis" class="space-y-4">
            <!-- Summary -->
            <div>
              <h4 class="font-semibold text-sm text-base-content/60 uppercase mb-1">Summary</h4>
              <p class="leading-relaxed">{{ parsedAnalysis.summary }}</p>
            </div>

            <!-- Sub-scores -->
            <div v-if="parsedAnalysis.scores" class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div v-for="item in scoreItems" :key="item.key" class="flex items-center gap-3">
                <div class="flex-1">
                  <div class="flex justify-between text-sm mb-1">
                    <span>{{ item.label }}</span>
                    <span class="font-mono font-bold">{{ item.value }}/{{ item.max }}</span>
                  </div>
                  <progress
                    class="progress w-full"
                    :class="progressClass(item.value, item.max)"
                    :value="item.value"
                    :max="item.max"
                  ></progress>
                </div>
              </div>
            </div>

            <!-- Strengths -->
            <div v-if="parsedAnalysis.strengths?.length">
              <h4 class="font-semibold text-sm text-success uppercase mb-1">Strengths</h4>
              <ul class="list-disc list-inside space-y-1">
                <li v-for="s in parsedAnalysis.strengths" :key="s">{{ s }}</li>
              </ul>
            </div>

            <!-- Weaknesses -->
            <div v-if="parsedAnalysis.weaknesses?.length">
              <h4 class="font-semibold text-sm text-error uppercase mb-1">Weaknesses</h4>
              <ul class="list-disc list-inside space-y-1">
                <li v-for="w in parsedAnalysis.weaknesses" :key="w">{{ w }}</li>
              </ul>
            </div>

            <!-- Recommendation -->
            <div v-if="parsedAnalysis.recommendation">
              <h4 class="font-semibold text-sm text-info uppercase mb-1">Recommendation</h4>
              <p class="leading-relaxed">{{ parsedAnalysis.recommendation }}</p>
            </div>
          </div>

          <!-- Failed -->
          <div v-else-if="data.analysis_status === 'failed'" class="flex items-center justify-between">
            <span class="text-error">Analysis failed</span>
            <button class="btn btn-sm btn-outline btn-error" :disabled="reanalyzeIsPending" @click="handleReanalyze">
              <span v-if="reanalyzeIsPending" class="loading loading-spinner loading-xs"></span>
              Retry
            </button>
          </div>

          <!-- Pending / Transcribing -->
          <div v-else class="flex items-center gap-3">
            <span class="loading loading-dots loading-sm"></span>
            <span class="text-base-content/50">Waiting for transcription to complete...</span>
          </div>
        </div>
      </div>

      <!-- Answers -->
      <div v-if="data.answers && data.answers.length > 0" class="card bg-base-100 shadow-md">
        <div class="card-body">
          <h3 class="card-title text-lg mb-4">Quiz Answers</h3>

          <div
            v-for="(qa, i) in data.answers"
            :key="qa.question_id"
            class="collapse collapse-arrow bg-base-200 mb-2"
          >
            <input type="radio" :name="'answers-accordion'" :checked="i === 0" />
            <div class="collapse-title font-medium flex items-center gap-2">
              <span class="badge badge-sm" :class="qa.score > 0 ? 'badge-success' : 'badge-error'">
                +{{ qa.score }}
              </span>
              Q{{ i + 1 }}: {{ qa.question_text }}
            </div>
            <div class="collapse-content">
              <div class="space-y-1 pt-2">
                <div
                  v-for="answer in qa.all_answers"
                  :key="answer.id"
                  class="flex items-center gap-2 p-2 rounded"
                  :class="answerRowClass(answer, qa.selected_answer_id)"
                >
                  <span
                    class="badge badge-xs shrink-0"
                    :class="answerBadgeClass(answer, qa.selected_answer_id)"
                  >
                    {{ answerIcon(answer, qa.selected_answer_id) }}
                  </span>
                  <span class="flex-1 min-w-0" :class="{ 'font-medium': answer.id === qa.selected_answer_id }">
                    {{ answer.text }}
                  </span>
                  <span class="badge badge-ghost badge-sm shrink-0">{{ answer.score }} pts</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </AdminLayout>
</template>

<script setup>
import { computed, watch, onUnmounted } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import AdminLayout from '../../components/admin/AdminLayout.vue'
import { useCandidateDetail, candidateAudioUrl, useReanalyzeCandidate } from '../../composables/useAdmin.js'

const route = useRoute()
const candidateId = computed(() => route.params.id)
const { data, isPending, error, refresh } = useCandidateDetail(candidateId)

const { mutate: reanalyze, isPending: reanalyzeIsPending } = useReanalyzeCandidate()

// Polling: refresh every 3s while analysis is in progress
let pollInterval = null

function startPolling() {
  stopPolling()
  pollInterval = setInterval(() => {
    refresh()
  }, 3000)
}

function stopPolling() {
  if (pollInterval) {
    clearInterval(pollInterval)
    pollInterval = null
  }
}

const isInProgress = computed(() => {
  if (!data.value) return false
  const status = data.value.analysis_status
  return status === 'pending' || status === 'transcribing' || status === 'analyzing'
})

const isTranscribing = computed(() => {
  if (!data.value) return false
  return data.value.analysis_status === 'transcribing'
})

const isAnalyzing = computed(() => {
  if (!data.value) return false
  return data.value.analysis_status === 'analyzing'
})

watch(isInProgress, (val) => {
  if (val) {
    startPolling()
  } else {
    stopPolling()
  }
}, { immediate: true })

onUnmounted(() => {
  stopPolling()
})

// Parse AI analysis JSON
const parsedAnalysis = computed(() => {
  if (!data.value?.ai_analysis) return null
  try {
    return JSON.parse(data.value.ai_analysis)
  } catch {
    return null
  }
})

const scoreItems = computed(() => {
  if (!parsedAnalysis.value?.scores) return []
  const s = parsedAnalysis.value.scores
  return [
    { key: 'descriptive_skills', label: 'Descriptive Skills', value: s.descriptive_skills || 0, max: 25 },
    { key: 'critical_thinking', label: 'Critical Thinking', value: s.critical_thinking || 0, max: 25 },
    { key: 'engagement', label: 'Engagement', value: s.engagement || 0, max: 20 },
    { key: 'structure', label: 'Structure', value: s.structure || 0, max: 15 },
    { key: 'practical_info', label: 'Practical Info', value: s.practical_info || 0, max: 15 },
  ]
})

function progressClass(value, max) {
  const pct = (value / max) * 100
  if (pct >= 70) return 'progress-success'
  if (pct >= 40) return 'progress-warning'
  return 'progress-error'
}

function handleReanalyze() {
  reanalyze(data.value.id)
  // Start polling after a short delay to give time for status update
  setTimeout(() => refresh(), 500)
}

const audioSrc = computed(() => {
  if (!data.value?.audio_path) return null
  const token = localStorage.getItem('admin_token')
  const url = candidateAudioUrl(data.value.id)
  return token ? `${url}?token=${token}` : url
})

const scoreColor = computed(() => {
  if (!data.value) return ''
  const s = data.value.total_score
  if (s >= 80) return 'text-success'
  if (s >= 50) return 'text-warning'
  return 'text-error'
})

const aiScoreBadgeClass = computed(() => {
  if (!data.value?.ai_score) return 'badge-ghost'
  const s = data.value.ai_score
  if (s >= 70) return 'badge-success'
  if (s >= 40) return 'badge-warning'
  return 'badge-error'
})

function answerRowClass(answer, selectedId) {
  if (answer.id === selectedId) return 'bg-primary/10'
  return ''
}

function answerBadgeClass(answer, selectedId) {
  if (answer.id === selectedId) return 'badge-primary'
  return 'badge-ghost'
}

function answerIcon(answer, selectedId) {
  if (answer.id === selectedId) return '\u2713'
  return '\u00A0'
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
</script>
