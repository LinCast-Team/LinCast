<template>
  <div class="flex flex-row px-2 py-1 mx-2 my-1">
      <div v-html="icon"></div>
      <h4 class="text-gray-400 text-lg">{{ title }}</h4>
  </div>
</template>

<script>
import feather from 'feather-icons';
import {
  computed,
} from 'vue';

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
  setup(props) {
    const icon = computed(() => feather.icons[props.iconClass].toSvg({ 'stroke-width': 1.5, class: 'text-gray-400 w-8 h-8 mx-4' }));

    return {
      icon,
    };
  },
};
</script>

<style>
</style>
