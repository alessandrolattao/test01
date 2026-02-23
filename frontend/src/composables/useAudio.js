import { ref, onUnmounted } from 'vue'

export function useAudioRecorder(maxDurationSec = 120) {
  const isRecording = ref(false)
  const audioBlob = ref(null)
  const audioUrl = ref(null)
  const timeLeft = ref(maxDurationSec)
  const error = ref(null)

  let mediaRecorder = null
  let chunks = []
  let timerInterval = null

  function startTimer() {
    timeLeft.value = maxDurationSec
    timerInterval = setInterval(() => {
      timeLeft.value--
      if (timeLeft.value <= 0) {
        stop()
      }
    }, 1000)
  }

  function stopTimer() {
    if (timerInterval) {
      clearInterval(timerInterval)
      timerInterval = null
    }
  }

  async function start() {
    try {
      error.value = null
      audioBlob.value = null
      audioUrl.value = null
      chunks = []

      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      mediaRecorder = new MediaRecorder(stream, { mimeType: 'audio/webm;codecs=opus' })

      mediaRecorder.ondataavailable = (e) => {
        if (e.data.size > 0) chunks.push(e.data)
      }

      mediaRecorder.onstop = () => {
        const blob = new Blob(chunks, { type: 'audio/webm' })
        audioBlob.value = blob
        audioUrl.value = URL.createObjectURL(blob)
        stream.getTracks().forEach((t) => t.stop())
        stopTimer()
      }

      mediaRecorder.start()
      isRecording.value = true
      startTimer()
    } catch (err) {
      error.value = 'Please allow microphone access to record your presentation.'
    }
  }

  function stop() {
    if (mediaRecorder && mediaRecorder.state === 'recording') {
      mediaRecorder.stop()
      isRecording.value = false
    }
  }

  function reset() {
    if (audioUrl.value) {
      URL.revokeObjectURL(audioUrl.value)
    }
    audioBlob.value = null
    audioUrl.value = null
    timeLeft.value = maxDurationSec
    error.value = null
  }

  function formatTime(seconds) {
    const m = Math.floor(seconds / 60)
    const s = seconds % 60
    return `${m}:${String(s).padStart(2, '0')}`
  }

  onUnmounted(() => {
    stopTimer()
    if (audioUrl.value) URL.revokeObjectURL(audioUrl.value)
  })

  return {
    isRecording,
    audioBlob,
    audioUrl,
    timeLeft,
    error,
    start,
    stop,
    reset,
    formatTime,
  }
}
