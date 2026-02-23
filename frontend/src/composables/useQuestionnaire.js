import { useQuery } from '@pinia/colada'
import { get } from '../api/client.js'

function getBrowserLang() {
  return navigator.language.split('-')[0]
}

export function useActiveQuestionnaire() {
  const lang = getBrowserLang()
  return useQuery({
    key: ['questionnaire', 'active', lang],
    query: () => get(`/questionnaire?lang=${lang}`),
  })
}
