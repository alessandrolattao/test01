<template>
  <div class="min-h-screen flex flex-col items-center justify-center p-4">
    <div class="w-full max-w-lg">
      <StepIndicator :current-step="1" />

      <div class="card bg-base-100 shadow-xl w-full">
        <div class="card-body">
          <h2 class="card-title text-2xl">{{ $t('candidate.home.title') }}</h2>
          <p class="text-base-content/70">{{ $t('candidate.home.subtitle') }}</p>

          <div v-if="apiError" class="alert alert-error mt-4">
            <span>{{ apiError }}</span>
          </div>

          <form class="mt-4" @submit.prevent="handleSubmit">
            <fieldset class="fieldset">
              <legend class="fieldset-legend">{{ $t('candidate.home.firstName') }}</legend>
              <input
                v-model="form.first_name"
                type="text"
                class="input input-bordered w-full"
                :placeholder="$t('candidate.home.firstNamePlaceholder')"
                required
              />
            </fieldset>

            <fieldset class="fieldset">
              <legend class="fieldset-legend">{{ $t('candidate.home.lastName') }}</legend>
              <input
                v-model="form.last_name"
                type="text"
                class="input input-bordered w-full"
                :placeholder="$t('candidate.home.lastNamePlaceholder')"
                required
              />
            </fieldset>

            <fieldset class="fieldset">
              <legend class="fieldset-legend">{{ $t('candidate.home.email') }}</legend>
              <input
                v-model="form.email"
                type="email"
                class="input input-bordered w-full"
                :placeholder="$t('candidate.home.emailPlaceholder')"
                required
              />
            </fieldset>

            <button
              type="submit"
              class="btn btn-primary btn-block mt-4"
              :disabled="asyncStatus === 'loading'"
            >
              <span v-if="asyncStatus === 'loading'" class="loading loading-spinner loading-sm"></span>
              {{ $t('candidate.home.submit') }}
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
import { useI18n } from 'vue-i18n'
import StepIndicator from '../components/candidate/StepIndicator.vue'
import { useRegisterCandidate } from '../composables/useCandidate.js'

const { t } = useI18n()
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
    apiError.value = err.message || t('candidate.home.error')
  }
}
</script>
