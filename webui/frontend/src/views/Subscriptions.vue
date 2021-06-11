<template>
  <div class="font-sans flex flex-col">
    <lc-header
      title='Subscriptions'
      class="mt-6"
    />
    <SubTabs />
  </div>
</template>

<script lang='ts'>
import LCHeader from '@/components/subscriptions/Header.vue';
// import ListItem from '@/components/library/ListItem.vue';
// import PodcastItem from '@/components/library/PodcastItem.vue';
import SubTabs from '@/components/subscriptions/SubTabs.vue';
import { ref } from 'vue';
// import feather from 'feather-icons';

export default {
  components: {
    'lc-header': LCHeader,
    // ListItem,
    // PodcastItem,
    SubTabs,
  },
  setup() {
    const podcastItem = ref([
      { title: 'Sample Podcast 1', key: 11 },
      { title: 'Sample Podcast 2', key: 12 },
      { title: 'Sample Podcast 3', key: 13 },
      { title: 'Sample Podcast 4', key: 14 },
      { title: 'Sample Podcast 5', key: 15 },
      { title: 'Sample Podcast 6', key: 16 },
    ]);
    return {
      // playIcon,
      podcastItem,
    };
  },
};
</script>

<style lang='scss' scoped>
@import "@/assets/css/_palette.scss";

.btn-bg {
  background-color: $primary-accent;
}
</style>
