<template>
  <div class="flex flex-row justify-between m-2 p-2">
    <div v-html="backIcon"></div>
    <h2 class="text-primary-dt text-xl font-semibold">{{ title }}</h2>
    <div v-html="moreIcon"></div>
  </div>
</template>

<script lang='ts'>
import { computed } from 'vue';
import feather from 'feather-icons';

export default {
  props: {
    title: {
      type: String,
      required: true,
    },
  },
  setup() {
    const backIcon = computed(() => feather.icons['arrow-left'].toSvg({ 'stroke-width': 1.5, class: 'text-secondary-dt w-7 h-7 mx-4' }));
    const moreIcon = computed(() => feather.icons['more-vertical'].toSvg({ 'stroke-width': 1.5, class: 'text-secondary-dt w-7 h-7 mx-4' }));

    return {
      backIcon,
      moreIcon,
    };
  },
};
</script>

<style lang='scss' scoped>
</style>
