<template>
  <div class="muted flex flex-col items-center m-2">
    <div v-html="toolIcon"></div>
    <p>{{ text }}</p>
  </div>
</template>

<script lang='ts'>
import { computed } from 'vue';
import feather from 'feather-icons';

export default {
  props: {
    text: {
      type: String,
      required: false,
      default: 'Work in progress',
    },
  },
  setup() {
    const toolIcon = computed(() => feather.icons.tool.toSvg({ 'stroke-width': 1.5, class: 'w-24 h-24' }));

    return {
      toolIcon,
    };
  },
};
</script>

<style lang="scss" scoped>
@import "@/assets/css/_palette.scss";

.muted {
  color: rgba($color: $text-secondary, $alpha: 0.3)
}
</style>
