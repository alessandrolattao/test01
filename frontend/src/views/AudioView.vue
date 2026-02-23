<template>
  <div class="min-h-screen flex flex-col items-center justify-center p-4">
    <div class="w-full max-w-2xl">
      <StepIndicator :current-step="3" />

      <!-- Upload error -->
      <div v-if="uploadError" class="alert alert-error mb-4">
        <span>{{ uploadError }}</span>
      </div>

      <AudioRecorderComponent
        :recorder="recorder"
        :submitting="uploadStatus === 'loading'"
        @submit="handleUpload"
      />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import StepIndicator from '../components/candidate/StepIndicator.vue'
import AudioRecorderComponent from '../components/candidate/AudioRecorder.vue'
import { useAudioRecorder } from '../composables/useAudio.js'
import { useUploadAudio } from '../composables/useCandidate.js'

const route = useRoute()
const router = useRouter()
const recorder = useAudioRecorder()
const { mutateAsync: uploadAudio, asyncStatus: uploadStatus } = useUploadAudio()
const uploadError = ref(null)

async function handleUpload() {
  try {
    uploadError.value = null
    await uploadAudio({
      candidateId: route.params.candidateId,
      audioBlob: recorder.audioBlob.value,
    })
    router.push({ name: 'done' })
  } catch (err) {
    uploadError.value = err.message || 'Failed to upload audio. Please try again.'
  }
}
</script>
