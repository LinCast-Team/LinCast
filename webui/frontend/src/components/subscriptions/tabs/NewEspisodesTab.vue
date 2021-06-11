<template>
  <h2 class="p-2 pt-3 text-secondary-dt">Today</h2>
  <episode-item
    title = "Episode"
    author = "Author"
    update = "Episode"
    resume = "Officia consequat sit nostrud nulla aliquip non. Dolore consequat officia mollit id consectetur. Enim deserunt do laboris minim occaecat nulla labore ad enim. Non et id minim ex excepteur. Occaecat aliqua veniam excepteur tempor voluptate."
    duration = "10:23"
  />
</template>

<script>
import EpisodeItem from '@/components/library/EpisodeItem.vue';

export default {
  name: 'NewEpisodesTab',
  components: {
    EpisodeItem,
  },
};
</script>

<style scoped>

</style>
