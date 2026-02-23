<template>
  <ul class="steps steps-horizontal w-full mb-8">
    <li class="step" :class="{ 'step-primary': currentStep >= 1 }">Registration</li>
    <li class="step" :class="{ 'step-primary': currentStep >= 2 }">Quiz</li>
    <li class="step" :class="{ 'step-primary': currentStep >= 3 }">Audio</li>
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
