<template>
  <!-- Remove this once implemented the search functionality -->
  <div class="flex flex-col items-center">
    <search
      :class="{ 'mt-4': !searchMode, 'm-0': searchMode }"
      @search-input="onSearchInput"
      :searchMode="searchMode"
      @search-focus="onSearchFocus"
    />

    <work-signal class="self-center absolute top-1/4">
      <div class="flex flex-col items-center">
        <i>Meanwhile, you can add the feed's URL directly there</i>
        <div v-html="smileIcon"></div>
      </div>
    </work-signal>

    <button class="btn px-10 py-1.5 rounded-2xl text-primary-dt absolute bottom-1/4" @click="submitFeed">
      Add
    </button>
  </div>

  <!-- Remove the 'hidden' class to show the content -->
  <div class="flex flex-col items-center justify-center font-sans hidden">
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
<script lang='ts'>
import { ref, computed } from 'vue';
import feather from 'feather-icons';
import Category from '@/components/discover/Category.vue';
import Search from '@/components/discover/Search.vue';
import Podcast from '@/components/discover/Podcast.vue';
import WorkSignal from '@/components/shared/WorkSignal.vue';
import { subscribe } from '@/store/api/subscriptions';

export default {
  components: {
    Category,
    Search,
    Podcast,
    WorkSignal,
  },
  setup() {
    const searchMode = ref(false);
    const content = ref('');

    const onSearchInput = (input: string) => {
      content.value = input;
    };

    const onSearchFocus = (f: boolean) => {
      searchMode.value = f;
    };

    const submitFeed = () => {
      if (content.value === '') {
        return;
      }

      subscribe(content.value)
        .then(() => {
          content.value = '';
        })
        .catch((err) => {
          console.error(err);
        });
    };

    const smileIcon = computed(() => feather.icons.smile.toSvg({ 'stroke-width': 2.0, class: 'w-6 h-6' }));

    return {
      onSearchInput,
      onSearchFocus,
      submitFeed,
      searchMode,
      smileIcon,
    };
  },
};
</script>
<style lang="scss">
@import '@/assets/css/_palette.scss';

.btn {
  background-color: $primary-accent;
}
</style>
