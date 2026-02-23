import { useQuery, useMutation, useQueryCache } from '@pinia/colada'
import { get, post } from '../api/client.js'

export function useAdminLogin() {
  return useMutation({
    mutation: (credentials) => post('/admin/login', credentials),
  })
}

export function useCandidatesList() {
  return useQuery({
    key: ['admin', 'candidates'],
    query: () => get('/admin/candidates'),
  })
}

export function useCandidateDetail(idRef) {
  return useQuery({
    key: () => ['admin', 'candidates', idRef.value],
    query: () => get(`/admin/candidates/${idRef.value}`),
    enabled: () => !!idRef.value,
  })
}

export function useQuestionnairesList() {
  return useQuery({
    key: ['admin', 'questionnaires'],
    query: () => get('/admin/questionnaires'),
  })
}

export function useQuestionnaireDetail(idRef) {
  return useQuery({
    key: () => ['admin', 'questionnaires', idRef.value],
    query: () => get(`/admin/questionnaires/${idRef.value}`),
    enabled: () => !!idRef.value,
  })
}

export function useCreateQuestionnaire() {
  const queryCache = useQueryCache()
  return useMutation({
    mutation: (data) => post('/admin/questionnaires', data),
    onSettled() {
      queryCache.invalidateQueries({ key: ['admin', 'questionnaires'] })
    },
  })
}

export function candidateAudioUrl(candidateId) {
  return `/api/admin/candidates/${candidateId}/audio`
}
