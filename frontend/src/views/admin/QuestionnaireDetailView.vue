<template>
  <AdminLayout>
    <!-- Loading -->
    <div v-if="isPending" class="flex justify-center items-center min-h-[50vh]">
      <span class="loading loading-spinner loading-lg"></span>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="alert alert-error">
      <span>{{ $t('admin.questionnaires.detailLoadError') }}</span>
      <button class="btn btn-sm btn-ghost" @click="refresh()">{{ $t('common.retry') }}</button>
    </div>

    <template v-else-if="data">
      <!-- Breadcrumb -->
      <div class="breadcrumbs text-sm mb-6">
        <ul>
          <li><RouterLink to="/admin/questionnaires">{{ $t('admin.questionnaires.title') }}</RouterLink></li>
          <li>{{ $t('admin.questionnaires.version', { v: data.version }) }}</li>
        </ul>
      </div>

      <!-- Header -->
      <div class="flex items-center gap-3 mb-2">
        <h1 class="text-2xl font-bold">{{ $t('admin.questionnaires.version', { v: data.version }) }}</h1>
        <div class="badge badge-lg" :class="data.is_active ? 'badge-success' : 'badge-ghost'">
          {{ data.is_active ? $t('admin.questionnaires.active') : $t('admin.questionnaires.inactive') }}
        </div>
      </div>
      <p class="text-base-content/50 mb-6">
        {{ $t('admin.questionnaires.createdOn', { date: formatDate(data.created_at), count: data.questions?.length || 0 }) }}
      </p>

      <!-- Questions -->
      <div
        v-for="(question, i) in data.questions"
        :key="question.id"
        class="card bg-base-100 shadow-md mb-4"
      >
        <div class="card-body">
          <h3 class="font-bold text-lg mb-3">
            <span class="badge badge-primary badge-outline mr-2">{{ i + 1 }}</span>
            {{ question.text }}
          </h3>
          <div class="space-y-2">
            <div
              v-for="answer in question.answers"
              :key="answer.id"
              class="flex items-center gap-3 p-2 rounded"
            >
              <span
                class="badge badge-sm"
                :class="answer.score > 0 ? 'badge-success' : 'badge-ghost'"
              >
                {{ answer.score }} {{ $t('common.pts') }}
              </span>
              <span :class="{ 'opacity-70': answer.score === 0 }">{{ answer.text }}</span>
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
import { useQuestionnaireDetail } from '../../composables/useAdmin.js'

const route = useRoute()
const questionnaireId = computed(() => route.params.id)
const { data, isPending, error, refresh } = useQuestionnaireDetail(questionnaireId)

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
</script>
