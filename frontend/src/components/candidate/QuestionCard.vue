<template>
  <div class="card bg-base-100 shadow-md mb-4 overflow-hidden">
    <div class="card-body">
      <div class="flex items-start gap-3">
        <span class="badge badge-primary badge-outline shrink-0 mt-1">{{ index + 1 }}</span>
        <h3 class="text-lg font-semibold leading-snug">{{ question.text }}</h3>
      </div>
      <div class="flex flex-col gap-2 mt-3">
        <label
          v-for="answer in question.answers"
          :key="answer.id"
          class="flex items-start gap-3 cursor-pointer p-3 rounded-lg hover:bg-base-200 transition-colors"
          :class="{ 'bg-primary/10': modelValue === answer.id }"
        >
          <input
            type="radio"
            :name="`question-${question.id}`"
            class="radio radio-primary shrink-0 mt-0.5"
            :value="answer.id"
            :checked="modelValue === answer.id"
            :disabled="disabled"
            @change="$emit('update:modelValue', answer.id)"
          />
          <span class="text-sm leading-relaxed">{{ answer.text }}</span>
        </label>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  question: { type: Object, required: true },
  index: { type: Number, required: true },
  modelValue: { type: String, default: null },
  disabled: { type: Boolean, default: false },
})

defineEmits(['update:modelValue'])
</script>
