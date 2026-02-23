<template>
  <div class="min-h-screen p-4">
    <div class="max-w-4xl mx-auto">
      <StepIndicator :current-step="2" />

      <!-- Loading -->
      <div v-if="isPending" class="flex justify-center p-12">
        <span class="loading loading-spinner loading-lg"></span>
      </div>

      <!-- Error -->
      <div v-else-if="error" class="alert alert-error">
        <span>Could not load the quiz. Please try again.</span>
        <button class="btn btn-sm btn-ghost" @click="refresh()">Retry</button>
      </div>

      <!-- Quiz content -->
      <template v-else-if="data">
        <!-- Progress -->
        <div class="mb-6">
          <div class="flex justify-between text-sm mb-1">
            <span>Progress</span>
            <span>{{ answeredCount }} of {{ data.questions.length }} answered</span>
          </div>
          <progress
            class="progress progress-primary w-full"
            :value="progressPct"
            max="100"
          ></progress>
        </div>

        <!-- Questions -->
        <QuestionCard
          v-for="(question, i) in data.questions"
          :key="question.id"
          :question="question"
          :index="i"
          :model-value="selectedAnswers[question.id] || null"
          :disabled="submitStatus === 'loading'"
          @update:model-value="(val) => (selectedAnswers[question.id] = val)"
        />

        <!-- Submit error -->
        <div v-if="submitError" class="alert alert-error mb-4">
          <span>{{ submitError }}</span>
        </div>

        <!-- Submit -->
        <div class="mt-6 mb-8">
          <button
            class="btn btn-primary btn-block btn-lg"
            :disabled="submitStatus === 'loading'"
            @click="handleSubmit"
          >
            <span v-if="submitStatus === 'loading'" class="loading loading-spinner loading-sm"></span>
            Submit Answers
          </button>
          <p v-if="!allAnswered" class="text-center text-sm text-base-content/50 mt-2">
            Answer all questions to continue
          </p>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { reactive, computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import StepIndicator from '../components/candidate/StepIndicator.vue'
import QuestionCard from '../components/candidate/QuestionCard.vue'
import { useActiveQuestionnaire } from '../composables/useQuestionnaire.js'
import { useSubmitAnswers } from '../composables/useCandidate.js'

const route = useRoute()
const router = useRouter()

const { data, isPending, error, refresh } = useActiveQuestionnaire()
const { mutateAsync: submitAnswers, asyncStatus: submitStatus } = useSubmitAnswers()

const selectedAnswers = reactive({})
const submitError = ref(null)

const answeredCount = computed(() => {
  if (!data.value) return 0
  return data.value.questions.filter((q) => selectedAnswers[q.id]).length
})

const allAnswered = computed(() => {
  if (!data.value) return false
  return answeredCount.value === data.value.questions.length
})

const progressPct = computed(() => {
  if (!data.value || data.value.questions.length === 0) return 0
  return Math.round((answeredCount.value / data.value.questions.length) * 100)
})

async function handleSubmit() {
  if (!allAnswered.value) return

  const answers = Object.entries(selectedAnswers).map(([questionId, answerId]) => ({
    question_id: questionId,
    answer_id: answerId,
  }))

  try {
    submitError.value = null
    await submitAnswers({ candidateId: route.params.candidateId, answers })
    router.push({ name: 'audio', params: { candidateId: route.params.candidateId } })
  } catch (err) {
    submitError.value = err.message || 'Failed to submit answers. Please try again.'
  }
}
</script>
