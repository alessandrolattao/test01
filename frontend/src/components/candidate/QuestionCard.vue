<template>
  <div class="card bg-base-100 shadow-md mb-4">
    <div class="card-body">
      <h3 class="card-title text-lg">
        <span class="badge badge-primary badge-outline mr-2">{{ index + 1 }}</span>
        {{ question.text }}
      </h3>
      <div class="flex flex-col gap-2 mt-3">
        <label
          v-for="answer in question.answers"
          :key="answer.id"
          class="label cursor-pointer justify-start gap-3 p-3 rounded-lg hover:bg-base-200 transition-colors"
          :class="{ 'bg-primary/10': modelValue === answer.id }"
        >
          <input
            type="radio"
            :name="`question-${question.id}`"
            class="radio radio-primary"
            :value="answer.id"
            :checked="modelValue === answer.id"
            :disabled="disabled"
            @change="$emit('update:modelValue', answer.id)"
          />
          <span class="label-text">{{ answer.text }}</span>
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
