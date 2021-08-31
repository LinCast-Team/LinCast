<template>
  <div class="flex p-2 py-5 justify-between">
    <podcast-item
      title= "Podcast"
      author= "Author"
    />
    <podcast-item
      title= "Podcast"
      author= "Author"
    />
    <podcast-item
      title= "Podcast"
      author= "Author"
    />
  </div>
</template>

<script>
import PodcastItem from '@/components/library/PodcastItem.vue';

export default {
  name: 'PodcastsTab',
  components: {
    'podcast-item': PodcastItem,
  },
};
</script>

<style scoped>

</style>
