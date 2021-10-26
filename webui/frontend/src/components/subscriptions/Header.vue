<template>
  <div class="flex flex-row justify-between py-2 px-4">
    <span @click="$router.go(-1)" v-html="backIcon"></span>
    <h2 class="text-primary-dt text-xl font-semibold">{{ title }}</h2>
    <span v-html="moreIcon"></span>
  </div>
</template>

<script lang='ts'>
import { defineComponent, computed } from 'vue';
import feather from 'feather-icons';

export default defineComponent({
  props: {
    title: {
      type: String,
      required: true,
    },
    route: {
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
});
</script>

<style lang='scss' scoped>
</style>
