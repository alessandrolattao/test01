import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { PiniaColada } from '@pinia/colada'
import router from './router'
import i18n from './i18n'
import App from './App.vue'
import './style.css'

const app = createApp(App)

const pinia = createPinia()
app.use(pinia)
app.use(PiniaColada)
app.use(router)
app.use(i18n)

app.mount('#app')
