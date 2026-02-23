import { createI18n } from 'vue-i18n'
import en from './en.json'
import it from './it.json'

const browserLang = navigator.language.split('-')[0]
const supportedLocales = ['en', 'it']
const locale = supportedLocales.includes(browserLang) ? browserLang : 'en'

const i18n = createI18n({
  legacy: false,
  locale,
  fallbackLocale: 'en',
  messages: { en, it },
})

export default i18n
