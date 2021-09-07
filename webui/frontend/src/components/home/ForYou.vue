<template>
  <div class="flex flex-col">
    <h1 class="text-primary-dt text-xl font-semibold">For You</h1>
    <work-signal class="self-center my-5"></work-signal>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import WorkSignal from '@/components/shared/WorkSignal.vue';

export default defineComponent({
  components: {
    WorkSignal,
  },
});
</script>
