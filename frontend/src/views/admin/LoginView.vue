<template>
  <div class="min-h-screen flex items-center justify-center bg-base-200 p-4">
    <div class="card bg-base-100 shadow-xl w-full max-w-sm">
      <div class="card-body">
        <h2 class="card-title text-2xl mb-4">{{ $t('admin.login.title') }}</h2>

        <div v-if="loginError" class="alert alert-error mb-4">
          <span>{{ loginError }}</span>
        </div>

        <form @submit.prevent="handleLogin">
          <fieldset class="fieldset">
            <legend class="fieldset-legend">{{ $t('admin.login.email') }}</legend>
            <input
              v-model="form.email"
              type="email"
              class="input input-bordered w-full"
              :placeholder="$t('admin.login.emailPlaceholder')"
              required
            />
          </fieldset>

          <fieldset class="fieldset">
            <legend class="fieldset-legend">{{ $t('admin.login.password') }}</legend>
            <input
              v-model="form.password"
              type="password"
              class="input input-bordered w-full"
              :placeholder="$t('admin.login.passwordPlaceholder')"
              required
            />
          </fieldset>

          <button
            type="submit"
            class="btn btn-primary btn-block mt-4"
            :disabled="asyncStatus === 'loading'"
          >
            <span v-if="asyncStatus === 'loading'" class="loading loading-spinner loading-sm"></span>
            {{ $t('admin.login.submit') }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAdminLogin } from '../../composables/useAdmin.js'

const { t } = useI18n()
const router = useRouter()
const { mutateAsync, asyncStatus } = useAdminLogin()
const loginError = ref(null)

const form = reactive({
  email: 'admin@test.com',
  password: 'password',
})

async function handleLogin() {
  try {
    loginError.value = null
    const result = await mutateAsync(form)
    localStorage.setItem('admin_token', result.token)
    router.push({ name: 'admin-candidates' })
  } catch (err) {
    loginError.value = err.message || t('admin.login.error')
  }
}
</script>
