import { useQuery } from '@pinia/colada'
import { get } from '../api/client.js'

export function useActiveQuestionnaire() {
  return useQuery({
    key: ['questionnaire', 'active'],
    query: () => get('/questionnaire'),
  })
}
