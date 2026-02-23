<template>
  <ul class="steps steps-horizontal w-full mb-8">
    <li class="step" :class="{ 'step-primary': currentStep >= 1 }">{{ $t('candidate.steps.registration') }}</li>
    <li class="step" :class="{ 'step-primary': currentStep >= 2 }">{{ $t('candidate.steps.quiz') }}</li>
    <li class="step" :class="{ 'step-primary': currentStep >= 3 }">{{ $t('candidate.steps.audio') }}</li>
  </ul>
</template>

<script setup>
defineProps({
  currentStep: {
    type: Number,
    required: true,
  },
})
</script>
