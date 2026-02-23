<template>
  <div class="card bg-base-100 shadow-xl">
    <div class="card-body items-center text-center">
      <h2 class="card-title text-2xl mb-2">Audio Presentation</h2>
      <p class="text-base-content/70 mb-6 max-w-md">
        Record a brief presentation about yourself.
        You have up to 2 minutes. Click the button below to start.
      </p>

      <!-- Error -->
      <div v-if="recorder.error.value" class="alert alert-warning mb-4">
        <span>{{ recorder.error.value }}</span>
      </div>

      <!-- Timer -->
      <div
        class="text-5xl font-mono font-bold mb-6 tabular-nums"
        :class="{ 'text-error': recorder.timeLeft.value < 30 && recorder.isRecording.value }"
      >
        {{ recorder.formatTime(recorder.timeLeft.value) }}
      </div>

      <!-- Recording indicator -->
      <div v-if="recorder.isRecording.value" class="flex items-center gap-2 mb-4">
        <span class="badge badge-error animate-pulse">REC</span>
        <span class="text-sm">Recording in progress...</span>
      </div>

      <!-- Record / Stop button (before recording exists) -->
      <template v-if="!recorder.audioUrl.value">
        <button
          class="btn btn-circle btn-lg w-20 h-20 mb-4"
          :class="recorder.isRecording.value ? 'btn-error' : 'btn-primary'"
          @click="recorder.isRecording.value ? recorder.stop() : recorder.start()"
        >
          <!-- Mic icon -->
          <svg v-if="!recorder.isRecording.value" xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-14 0m7 7v4m-4 0h8m-4-18a3 3 0 00-3 3v4a3 3 0 006 0V7a3 3 0 00-3-3z" />
          </svg>
          <!-- Stop icon -->
          <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="currentColor" viewBox="0 0 24 24">
            <rect x="6" y="6" width="12" height="12" rx="2" />
          </svg>
        </button>
        <p class="text-sm text-base-content/50">
          {{ recorder.isRecording.value ? 'Click to stop recording' : 'Click to start recording' }}
        </p>
      </template>

      <!-- Playback (after recording) -->
      <template v-if="recorder.audioUrl.value && !recorder.isRecording.value">
        <div class="w-full space-y-4">
          <audio :src="recorder.audioUrl.value" controls class="w-full"></audio>
          <div class="flex gap-3">
            <button class="btn btn-ghost flex-1" :disabled="submitting" @click="recorder.reset()">
              Record Again
            </button>
            <button class="btn btn-primary flex-1" :disabled="submitting" @click="$emit('submit')">
              <span v-if="submitting" class="loading loading-spinner loading-sm"></span>
              {{ submitting ? 'Uploading...' : 'Submit Recording' }}
            </button>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
defineProps({
  recorder: { type: Object, required: true },
  submitting: { type: Boolean, default: false },
})

defineEmits(['submit'])
</script>
