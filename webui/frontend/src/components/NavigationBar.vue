<template>
  <div class="grid grid-cols-3 z-50 justify-items-center p-4 transition-colors duration-500 text-gray-50 bg-gradient-to-br from-gray-700 to-gray-900 ">
    <div v-html="homeIcon"></div>
    <div v-html="searchIcon"></div>
    <div v-html="bookmarkIcon"></div>
  </div>
</template>

<script>
import { computed } from 'vue';
import feather from 'feather-icons';

export default {
  setup() {
    const homeIcon = computed(() => feather.icons.home.toSvg()); /* { 'stroke-width': 0.8, class: 'w-16 h-16 md:w-20 md:h-20' } */
    const searchIcon = computed(() => feather.icons.search.toSvg());
    const bookmarkIcon = computed(() => feather.icons.bookmark.toSvg());

    return {
      homeIcon,
      searchIcon,
      bookmarkIcon,
    };
  },
};
</script>

<style>
</style>
