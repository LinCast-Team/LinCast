<template>
<div class="flex px-2 py-5">
  <div v-if="episodes.length > 0">
    <h2 class="p-2 pt-3 text-secondary-dt">Today</h2>
    <!-- FIX Some props do not match with the fields of the object -->
    <episode-item
      v-for="e in episodes"
      :title="e.title"
      :author="e.podcastTitle"
      :resume="e.description"
      :duration="e.duration"
      :key="e.id"
    />
    <!-- <episode-item
      title = "Episode"
      author = "Author"
      resume = "Officia consequat sit nostrud nulla aliquip non. Dolore consequat officia mollit id consectetur. Enim deserunt do laboris minim occaecat nulla labore ad enim. Non et id minim ex excepteur. Occaecat aliqua veniam excepteur tempor voluptate."
      duration = "10:23"
    /> -->
  </div>
  <div v-else class="mx-auto">
    <p class="self-center text-secondary-dt">Looks like there is nothing here...</p>
  </div>
</div>
</template>

<script lang='ts'>
import { defineComponent, ref } from 'vue';
import EpisodeItem from '@/components/library/EpisodeItem.vue';

export default defineComponent({
  components: {
    EpisodeItem,
  },
  setup() {
    const episodes = ref([]);

    return {
      episodes,
    };
  },
});
</script>

<style scoped>

</style>
