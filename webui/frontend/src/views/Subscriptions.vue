<template>
  <div class="font-sans flex flex-col">
    <lc-header
      title='Subscriptions'
      class="mt-6"
    />
    <SubTabs />
    <!-- <div class="flex justify-evenly border-b-2 border-gray-700">
      <button class="primary-accent w-28 h-8 border-b-2 border-green-600 rounded-none">Podcasts</button>
      <button class="text-secondary-dt w-28 h-8">New Episodes</button>
    </div> -->
    <!-- <podcast-item
          v-for="rp in PodcastsItem"
          :title="rp.title"
          :key="rp.key"
          imgSrc='https://picsum.photos/600'
        /> -->
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
