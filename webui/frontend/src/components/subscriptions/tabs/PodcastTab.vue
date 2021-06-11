<template>
  <div class="flex p-2 py-5">
    <podcast-item
      title= "Podcast"
      update= "Day ago"
    />
    <podcast-item
      title= "Podcast"
      update= "Day ago"
    />
    <podcast-item
      title= "Podcast"
      update= "Day ago"
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
