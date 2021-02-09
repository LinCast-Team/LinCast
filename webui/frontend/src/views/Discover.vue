<template>
  <div class="flex flex-col items-center justify-center font-sans">
    <search
      :class="{ 'mt-4': !searchMode, 'm-0': searchMode }"
      @search-input="onSearchInput"
      :searchMode="searchMode"
      @search-focus="onSearchFocus"
    />

    <div v-show="!searchMode">
      <h3 class="my-4 w-64 text-lg text-gray-50 inline-block font-semibold">Categories</h3>

      <category title="Arts"/>
      <category title="Business"/>
      <category title="News"/>
      <category title="Music"/>
      <category title="Technology"/>
    </div>

    <div v-show="searchMode" class="w-full">
      <podcast title="Hello world 1" author="Martin"/>
      <podcast title="Hello world 2" author="Bruno"/>
      <podcast title="Hello world 3" author="Martin"/>
      <podcast title="Hello world 4" author="Bruno"/>
      <podcast title="He" author="Bruno"/>
    </div>
  </div>
</template>
<script>
import { ref } from 'vue';
import Category from '@/components/discover/Category.vue';
import Search from '@/components/discover/Search.vue';
import Podcast from '@/components/discover/Podcast.vue';

export default {
  components: {
    Category,
    Search,
    Podcast,
  },
  setup() {
    const searchMode = ref(false);

    const onSearchInput = (input) => {
      console.log(input);
    };

    const onSearchFocus = (f) => {
      searchMode.value = f;
    };

    return {
      onSearchInput,
      onSearchFocus,
      searchMode,
    };
  },
};
</script>
<style lang="scss">
</style>
