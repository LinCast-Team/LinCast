<template>
  <div class="flex flex-col">
    <div class="flex flex-row justify-between">
      <h2 class="text-primary-dt text-xl font-semibold">Subscriptions</h2>
      <router-link
        :to="{ name: 'Subscriptions' }"
        class="link text-sm self-center"
        >Show more</router-link>
    </div>
    <work-signal class="self-center my-5"/>
  </div>
</template>

<script lang='ts'>
import { defineComponent } from 'vue';
import WorkSignal from '@/components/shared/WorkSignal.vue';

export default defineComponent({
  components: {
    WorkSignal,
  },
});
</script>

<style lang='scss' scoped>
</style>
