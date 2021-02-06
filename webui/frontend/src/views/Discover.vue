<template>
  <div class="flex flex-col items-center justify-center font-sans">
    <search class="mt-20" @search-input="onSearchInput"/>

    <h3 class="my-4 w-64 text-lg text-gray-50 inline-block font-semibold">Categories</h3>

    <category title="Arts"/>
    <category title="Business"/>
    <category title="News"/>
    <category title="Music"/>
    <category title="Technology"/>
  </div>
</template>
<script>
import Category from '@/components/discover/Category.vue';
import Search from '@/components/discover/Search.vue';

export default {
  components: {
    Category,
    Search,
  },
  setup() {
    const onSearchInput = (input) => {
      console.log(input);
    };

    return {
      onSearchInput,
    };
  },
};
</script>
<style lang="scss">
</style>
