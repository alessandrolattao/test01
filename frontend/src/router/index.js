import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // Candidate flow
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomeView.vue'),
    },
    {
      path: '/quiz/:candidateId',
      name: 'quiz',
      component: () => import('../views/QuizView.vue'),
    },
    {
      path: '/audio/:candidateId',
      name: 'audio',
      component: () => import('../views/AudioView.vue'),
    },
    {
      path: '/done',
      name: 'done',
      component: () => import('../views/DoneView.vue'),
    },
    // Admin flow
    {
      path: '/admin/login',
      name: 'admin-login',
      component: () => import('../views/admin/LoginView.vue'),
    },
    {
      path: '/admin/candidates',
      name: 'admin-candidates',
      component: () => import('../views/admin/CandidatesView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/admin/candidates/:id',
      name: 'admin-candidate-detail',
      component: () => import('../views/admin/CandidateDetailView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/admin/questionnaires',
      name: 'admin-questionnaires',
      component: () => import('../views/admin/QuestionnairesView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/admin/questionnaires/new',
      name: 'admin-questionnaire-new',
      component: () => import('../views/admin/QuestionnaireNewView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/admin/questionnaires/:id',
      name: 'admin-questionnaire-detail',
      component: () => import('../views/admin/QuestionnaireDetailView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !localStorage.getItem('admin_token')) {
    return { name: 'admin-login' }
  }
})

export default router
