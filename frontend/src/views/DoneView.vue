<template>
  <div class="min-h-screen flex items-center justify-center p-4">
    <div class="card bg-base-100 shadow-xl max-w-md w-full">
      <div class="card-body items-center text-center">
        <div class="w-20 h-20 rounded-full bg-success/20 flex items-center justify-center mb-4">
          <svg class="w-10 h-10 text-success" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h2 class="card-title text-2xl">{{ $t('candidate.done.title') }}</h2>
        <p class="text-base-content/70">{{ $t('candidate.done.message') }}</p>
      </div>
    </div>
  </div>
</template>
