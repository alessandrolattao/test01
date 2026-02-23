<template>
  <div class="drawer lg:drawer-open">
    <input id="admin-drawer" type="checkbox" class="drawer-toggle" />

    <!-- Main content area -->
    <div class="drawer-content">
      <!-- Top navbar (mobile only) -->
      <div class="navbar bg-base-200 lg:hidden">
        <label for="admin-drawer" class="btn btn-square btn-ghost">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </label>
        <span class="font-bold ml-2">{{ $t('admin.layout.title') }}</span>
      </div>

      <div class="p-4 md:p-8">
        <slot />
      </div>
    </div>

    <!-- Sidebar -->
    <div class="drawer-side">
      <label for="admin-drawer" class="drawer-overlay"></label>
      <div class="bg-base-200 min-h-full w-64 p-4 flex flex-col">
        <div class="font-bold text-lg mb-6 px-2">{{ $t('admin.layout.title') }}</div>

        <ul class="menu flex-1">
          <li>
            <RouterLink to="/admin/candidates" :class="{ active: isCandidatesActive }">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z" />
              </svg>
              {{ $t('admin.layout.candidates') }}
            </RouterLink>
          </li>
          <li>
            <RouterLink to="/admin/questionnaires" :class="{ active: isQuestionnairesActive }">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              {{ $t('admin.layout.questionnaires') }}
            </RouterLink>
          </li>
        </ul>

        <div class="mt-auto">
          <button class="btn btn-ghost btn-block justify-start" @click="logout">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
            {{ $t('admin.layout.logout') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const isCandidatesActive = computed(() => route.path.startsWith('/admin/candidates'))
const isQuestionnairesActive = computed(() => route.path.startsWith('/admin/questionnaires'))

function logout() {
  localStorage.removeItem('admin_token')
  router.push('/admin/login')
}
</script>
