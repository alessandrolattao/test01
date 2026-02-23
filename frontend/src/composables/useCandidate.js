import { useMutation } from '@pinia/colada'
import { post } from '../api/client.js'

export function useRegisterCandidate() {
  return useMutation({
    mutation: (data) => post('/candidates', data),
  })
}

export function useSubmitAnswers() {
  return useMutation({
    mutation: ({ candidateId, answers }) =>
      post(`/candidates/${candidateId}/answers`, { answers }),
  })
}

export function useUploadAudio() {
  return useMutation({
    mutation: ({ candidateId, audioBlob }) => {
      const formData = new FormData()
      formData.append('audio', audioBlob, 'recording.webm')
      return post(`/candidates/${candidateId}/audio`, formData)
    },
  })
}
