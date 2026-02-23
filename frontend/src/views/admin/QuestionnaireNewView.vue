<template>
  <AdminLayout>
    <!-- Breadcrumb -->
    <div class="breadcrumbs text-sm mb-6">
      <ul>
        <li><RouterLink to="/admin/questionnaires">Questionnaires</RouterLink></li>
        <li>New Questionnaire</li>
      </ul>
    </div>

    <h1 class="text-2xl font-bold mb-6">New Questionnaire</h1>

    <!-- API error -->
    <div v-if="saveError" class="alert alert-error mb-4">
      <span>{{ saveError }}</span>
    </div>

    <!-- Questions -->
    <div
      v-for="(question, qi) in questions"
      :key="qi"
      class="card bg-base-100 shadow-md mb-4"
    >
      <div class="card-body">
        <div class="flex items-center justify-between mb-3">
          <h3 class="font-bold text-lg">Question {{ qi + 1 }}</h3>
          <button
            class="btn btn-ghost btn-sm text-error"
            :disabled="questions.length <= 1"
            @click="removeQuestion(qi)"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            Remove
          </button>
        </div>

        <fieldset class="fieldset mb-4">
          <legend class="fieldset-legend">Question Text</legend>
          <input
            v-model="question.text"
            type="text"
            class="input input-bordered w-full"
            placeholder="Enter your question..."
          />
          <p v-if="question.textError" class="text-error text-sm mt-1">{{ question.textError }}</p>
        </fieldset>

        <div class="space-y-2">
          <p class="font-medium text-sm">Answers</p>

          <div
            v-for="(answer, ai) in question.answers"
            :key="ai"
            class="flex items-center gap-2"
          >
            <input
              v-model="answer.text"
              type="text"
              class="input input-bordered flex-1"
              placeholder="Answer text"
            />
            <input
              v-model.number="answer.score"
              type="number"
              class="input input-bordered w-20 text-center"
              placeholder="Score"
              min="0"
            />
            <button
              class="btn btn-ghost btn-sm btn-square text-error"
              :disabled="question.answers.length <= 2"
              @click="removeAnswer(qi, ai)"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <p v-if="question.answersError" class="text-error text-sm">{{ question.answersError }}</p>

          <button class="btn btn-ghost btn-sm text-secondary mt-2" @click="addAnswer(qi)">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Add Answer
          </button>
        </div>
      </div>
    </div>

    <!-- Add question -->
    <button class="btn btn-outline btn-secondary btn-block mb-6" @click="addQuestion">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      Add Question
    </button>

    <!-- Save / Cancel -->
    <div class="flex justify-end gap-3 mt-6 mb-8">
      <RouterLink to="/admin/questionnaires" class="btn btn-ghost">Cancel</RouterLink>
      <button
        class="btn btn-primary"
        :disabled="asyncStatus === 'loading'"
        @click="handleSave"
      >
        <span v-if="asyncStatus === 'loading'" class="loading loading-spinner loading-sm"></span>
        Save Questionnaire
      </button>
    </div>

    <!-- Toast -->
    <div v-if="showToast" class="toast toast-end">
      <div class="alert alert-success">
        <span>Questionnaire created successfully</span>
      </div>
    </div>
  </AdminLayout>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import AdminLayout from '../../components/admin/AdminLayout.vue'
import { useCreateQuestionnaire } from '../../composables/useAdmin.js'

const router = useRouter()
const { mutateAsync, asyncStatus } = useCreateQuestionnaire()
const saveError = ref(null)
const showToast = ref(false)

function makeAnswer() {
  return { text: '', score: 0 }
}

function makeQuestion() {
  return {
    text: '',
    textError: null,
    answersError: null,
    answers: [makeAnswer(), makeAnswer()],
  }
}

const questions = reactive([makeQuestion()])

function addQuestion() {
  questions.push(makeQuestion())
}

function removeQuestion(qi) {
  if (questions.length > 1) questions.splice(qi, 1)
}

function addAnswer(qi) {
  questions[qi].answers.push(makeAnswer())
}

function removeAnswer(qi, ai) {
  if (questions[qi].answers.length > 2) {
    questions[qi].answers.splice(ai, 1)
  }
}

function validate() {
  let valid = true
  for (const q of questions) {
    q.textError = null
    q.answersError = null

    if (!q.text.trim()) {
      q.textError = 'Question text is required'
      valid = false
    }

    const emptyAnswers = q.answers.some((a) => !a.text.trim())
    if (emptyAnswers) {
      q.answersError = 'All answers must have text'
      valid = false
    }

    const hasPositiveScore = q.answers.some((a) => a.score > 0)
    if (!hasPositiveScore) {
      q.answersError = (q.answersError ? q.answersError + '. ' : '') +
        'At least one answer must have a score greater than 0'
      valid = false
    }
  }
  return valid
}

async function handleSave() {
  if (!validate()) return

  const payload = {
    questions: questions.map((q, i) => ({
      text: q.text,
      sort_order: i + 1,
      answers: q.answers.map((a, j) => ({
        text: a.text,
        score: a.score,
        sort_order: j + 1,
      })),
    })),
  }

  try {
    saveError.value = null
    const result = await mutateAsync(payload)
    showToast.value = true
    setTimeout(() => {
      router.push({ name: 'admin-questionnaire-detail', params: { id: result.id } })
    }, 1000)
  } catch (err) {
    saveError.value = err.message || 'Failed to create questionnaire.'
  }
}
</script>
