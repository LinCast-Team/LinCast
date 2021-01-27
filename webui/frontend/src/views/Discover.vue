<template>
  <div class="flex flex-col items-center justify-center font-sans">
    <div>
      <div class="mt-20 p-0.5 rounded-md bg-gradient-to-r from-indigo-700 to-green-700">
        <div class="p-3 w-64 bg-gray-800 flex">
          <div v-html="searchIcon"></div>
          <input
            class="rounded-md bg-transparent border-0 placeholder-white focus:placeholder-gray-500 focus:outline-none"
            type="text"
            role="search"
            name="Search a podcast"
            placeholder="Search"
          >
        </div>
      </div>
    </div>
    <div class="m-2 p-0.5 ">
      <h3 class="p-2 w-64 text-lg text-green-600 inline-block font-semibold">Categories</h3>
    </div>
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
    const searchIcon = computed(() => feather.icons.search.toSvg({ 'stroke-width': 1.5, class: 'w-5 h-5 md:w-12 md:h-12 text-white mr-2' }));
    return {
      searchIcon,
    };
  },
};
</script>
<style lang="scss">
</style>
