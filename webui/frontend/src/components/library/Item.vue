<template>
  <div class="flex flex-row px-2 py-1 mx-2 my-3 text-secondary-dt">
      <div v-html="icon"></div>
      <h4 class="text-lg">{{ title }}</h4>
  </div>
</template>

<script lang='ts'>
import feather from 'feather-icons';
import {
  computed,
} from 'vue';

interface Data {
  [key: string]: unknown;
}

export default {
  props: {
    title: {
      type: String,
      required: true,
    },
    iconClass: {
      type: String,
      required: true,
    },
  },
  setup(props: Data): Data {
    const icon = computed(() => feather.icons[props.iconClass as string].toSvg({ 'stroke-width': 1.5, class: 'w-7 h-7 mx-4' }));

    return {
      icon,
    };
  },
};
</script>

<style>
</style>
