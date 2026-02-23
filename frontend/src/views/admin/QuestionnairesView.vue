<template>
  <AdminLayout>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">{{ $t('admin.questionnaires.title') }}</h1>
      <RouterLink to="/admin/questionnaires/new" class="btn btn-primary">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        {{ $t('admin.questionnaires.new') }}
      </RouterLink>
    </div>

    <!-- Loading -->
    <div v-if="isPending" class="flex justify-center p-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="alert alert-error">
      <span>{{ $t('admin.questionnaires.loadError') }}</span>
      <button class="btn btn-sm btn-ghost" @click="refresh()">{{ $t('common.retry') }}</button>
    </div>

    <!-- Empty -->
    <div v-else-if="data && data.length === 0" class="text-center py-12 text-base-content/50">
      <p class="text-lg mb-4">{{ $t('admin.questionnaires.empty') }}</p>
      <RouterLink to="/admin/questionnaires/new" class="btn btn-primary">
        {{ $t('admin.questionnaires.createFirst') }}
      </RouterLink>
    </div>

    <!-- List -->
    <div v-else-if="data" class="space-y-3">
      <div
        v-for="q in data"
        :key="q.id"
        class="card bg-base-100 shadow-md"
        :class="{ 'border-l-4 border-success': q.is_active }"
      >
        <div class="card-body flex-row items-center justify-between py-4">
          <div>
            <div class="flex items-center gap-2">
              <h3 class="font-bold">{{ $t('admin.questionnaires.version', { v: q.version }) }}</h3>
              <div class="badge" :class="q.is_active ? 'badge-success' : 'badge-ghost'">
                {{ q.is_active ? $t('admin.questionnaires.active') : $t('admin.questionnaires.inactive') }}
              </div>
            </div>
            <p class="text-sm text-base-content/50 mt-1">
              {{ $t('admin.questionnaires.questionsCount', { count: q.questions_count || 0, date: formatDate(q.created_at) }) }}
            </p>
          </div>
          <RouterLink
            :to="{ name: 'admin-questionnaire-detail', params: { id: q.id } }"
            class="btn btn-ghost btn-sm"
          >
            {{ $t('admin.questionnaires.view') }}
          </RouterLink>
        </div>
      </div>
    </div>
  </AdminLayout>
</template>

<script setup>
import { RouterLink } from 'vue-router'
import AdminLayout from '../../components/admin/AdminLayout.vue'
import { useQuestionnairesList } from '../../composables/useAdmin.js'

const { data, isPending, error, refresh } = useQuestionnairesList()

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
</script>
