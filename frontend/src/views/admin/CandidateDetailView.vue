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
                    class="badge badge-xs"
                    :class="answerBadgeClass(answer, qa.selected_answer_id)"
                  >
                    {{ answerIcon(answer, qa.selected_answer_id) }}
                  </span>
                  <span :class="{ 'font-medium': answer.id === qa.selected_answer_id }">
                    {{ answer.text }}
                  </span>
                  <span class="badge badge-ghost badge-sm ml-auto">{{ answer.score }} pts</span>
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
import { computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import AdminLayout from '../../components/admin/AdminLayout.vue'
import { useCandidateDetail, candidateAudioUrl } from '../../composables/useAdmin.js'

const route = useRoute()
const candidateId = computed(() => route.params.id)
const { data, isPending, error, refresh } = useCandidateDetail(candidateId)

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
