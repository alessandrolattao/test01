<template>
  <div class="min-h-screen flex flex-col items-center justify-center p-4">
    <div class="w-full max-w-lg">
      <StepIndicator :current-step="1" />

      <div class="card bg-base-100 shadow-xl w-full">
        <div class="card-body">
          <h2 class="card-title text-2xl">Join Our Team</h2>
          <p class="text-base-content/70">Fill in your details to get started</p>

          <div v-if="apiError" class="alert alert-error mt-4">
            <span>{{ apiError }}</span>
          </div>

          <form class="mt-4" @submit.prevent="handleSubmit">
            <fieldset class="fieldset">
              <legend class="fieldset-legend">First Name</legend>
              <input
                v-model="form.first_name"
                type="text"
                class="input input-bordered w-full"
                placeholder="Mario"
                required
              />
            </fieldset>

            <fieldset class="fieldset">
              <legend class="fieldset-legend">Last Name</legend>
              <input
                v-model="form.last_name"
                type="text"
                class="input input-bordered w-full"
                placeholder="Rossi"
                required
              />
            </fieldset>

            <fieldset class="fieldset">
              <legend class="fieldset-legend">Email</legend>
              <input
                v-model="form.email"
                type="email"
                class="input input-bordered w-full"
                placeholder="mario@example.com"
                required
              />
            </fieldset>

            <button
              type="submit"
              class="btn btn-primary btn-block mt-4"
              :disabled="asyncStatus === 'loading'"
            >
              <span v-if="asyncStatus === 'loading'" class="loading loading-spinner loading-sm"></span>
              Start Quiz
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import StepIndicator from '../components/candidate/StepIndicator.vue'
import { useRegisterCandidate } from '../composables/useCandidate.js'

const router = useRouter()
const apiError = ref(null)

const form = reactive({
  first_name: '',
  last_name: '',
  email: '',
})

const { mutateAsync, asyncStatus } = useRegisterCandidate()

async function handleSubmit() {
  try {
    apiError.value = null
    const result = await mutateAsync(form)
    router.push({ name: 'quiz', params: { candidateId: result.id } })
  } catch (err) {
    apiError.value = err.message || 'Registration failed. Please try again.'
  }
}
</script>
