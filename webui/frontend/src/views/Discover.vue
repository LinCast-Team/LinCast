<template>
  <div class="flex flex-col items-center justify-center font-sans">
    <div>
      <div class="mt-20 p-3 w-64 bg-gray-800 flex rounded-xl border border-green-700">
        <div v-html="searchIcon"></div>
        <input
          class="w-auto bg-transparent border-0 placeholder-gray-50 text-gray-50 focus:placeholder-gray-500 focus:outline-none"
          type="text"
          role="search"
          name="Search a podcast"
          placeholder="Search"
          aria-placeholder="Search"
        >
      </div>
    </div>

    <h3 class="my-4 w-64 text-lg text-gray-50 inline-block font-semibold">Categories</h3>

    <category title="Arts"/>
    <category title="Business"/>
    <category title="News"/>
    <category title="Music"/>
    <category title="Technology"/>
  </div>
</template>
<script>
import { computed } from 'vue';
import feather from 'feather-icons';
import Category from '@/components/discover/Category.vue';

export default {
  components: {
    Category,
  },
  setup() {
    const searchIcon = computed(() => feather.icons.search.toSvg({ 'stroke-width': 1.5, class: 'w-5 h-5 md:w-12 md:h-12 text-gray-50 mr-2' }));
    return {
      searchIcon,
    };
  },
};
</script>
<style lang="scss">
</style>
