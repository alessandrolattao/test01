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
        <!-- Progress bar -->
        <div class="mb-6">
          <div class="flex justify-between text-sm mb-1">
            <span>Question {{ currentIndex + 1 }} of {{ totalQuestions }}</span>
            <span>{{ answeredCount }} answered</span>
          </div>
          <progress
            class="progress progress-primary w-full"
            :value="progressPct"
            max="100"
          ></progress>
        </div>

        <!-- Single question -->
        <QuestionCard
          :question="currentQuestion"
          :index="currentIndex"
          :model-value="selectedAnswers[currentQuestion.id] || null"
          :disabled="submitStatus === 'loading'"
          @update:model-value="handleAnswer"
        />

        <!-- Submit error -->
        <div v-if="submitError" class="alert alert-error mb-4">
          <span>{{ submitError }}</span>
        </div>

        <!-- Navigation -->
        <div class="flex gap-3 mt-6 mb-8">
          <button
            v-if="currentIndex > 0"
            class="btn btn-outline flex-1"
            @click="currentIndex--"
          >
            Previous
          </button>
          <div v-else class="flex-1"></div>

          <button
            v-if="!isLastQuestion"
            class="btn btn-primary flex-1"
            :disabled="!currentAnswered"
            @click="currentIndex++"
          >
            Next
          </button>
          <button
            v-else
            class="btn btn-primary flex-1"
            :disabled="!allAnswered || submitStatus === 'loading'"
            @click="handleSubmit"
          >
            <span v-if="submitStatus === 'loading'" class="loading loading-spinner loading-sm"></span>
            Submit Answers
          </button>
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
const currentIndex = ref(0)

const totalQuestions = computed(() => data.value?.questions.length || 0)

const currentQuestion = computed(() => data.value?.questions[currentIndex.value])

const isLastQuestion = computed(() => currentIndex.value === totalQuestions.value - 1)

const currentAnswered = computed(() => !!selectedAnswers[currentQuestion.value?.id])

const answeredCount = computed(() => {
  if (!data.value) return 0
  return data.value.questions.filter((q) => selectedAnswers[q.id]).length
})

const allAnswered = computed(() => {
  if (!data.value) return false
  return answeredCount.value === totalQuestions.value
})

const progressPct = computed(() => {
  if (totalQuestions.value === 0) return 0
  return Math.round((answeredCount.value / totalQuestions.value) * 100)
})

function handleAnswer(answerId) {
  selectedAnswers[currentQuestion.value.id] = answerId
  // Auto-advance after a short delay
  if (!isLastQuestion.value) {
    setTimeout(() => {
      currentIndex.value++
    }, 300)
  }
}

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
