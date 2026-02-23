import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
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
    {
      path: '/admin/login',
      name: 'admin-login',
      component: () => import('../views/admin/LoginView.vue'),
    },
    {
      path: '/admin/candidates',
      name: 'admin-candidates',
      component: () => import('../views/admin/CandidatesView.vue'),
    },
    {
      path: '/admin/candidates/:id',
      name: 'admin-candidate-detail',
      component: () => import('../views/admin/CandidateDetailView.vue'),
    },
    {
      path: '/admin/questionnaires',
      name: 'admin-questionnaires',
      component: () => import('../views/admin/QuestionnairesView.vue'),
    },
    {
      path: '/admin/questionnaires/new',
      name: 'admin-questionnaire-new',
      component: () => import('../views/admin/QuestionnaireNewView.vue'),
    },
    {
      path: '/admin/questionnaires/:id',
      name: 'admin-questionnaire-detail',
      component: () => import('../views/admin/QuestionnaireDetailView.vue'),
    },
  ],
})

export default router
