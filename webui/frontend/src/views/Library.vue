<template>
  <div class="font-sans flex flex-col overflow-y-auto">
    <div class="mt-6">
      <h3 class="text-gray-50 my-2 mx-8 text-lg font-semibold">Latest</h3>
      <div class="flex flex-row my-2 overflow-x-auto">
        <recent-podcast
          v-for="rp in recentPodcasts"
          :title="rp.title"
          :key="rp.key"
          imgSrc='https://picsum.photos/600'
        />
      </div>
      <div class="w-full bg-gray-700" style="height: 1px"></div>
      <div class="p-4">
        <item title="Listen Later" icon-class="package" />
        <item title="Subscriptions" icon-class="layers" />
        <item title="Latest Episodes" icon-class="clock" />
        <item title="History" icon-class="list" />
      </div>
      <div class="w-full bg-gray-700" style="height: 1px"></div>
    </div>
    <lincast-info />
  </div>
</template>

<script lang='ts'>
import Item from '@/components/library/Item.vue';
import RecentPodcast from '@/components/library/RecentPodcast.vue';
import LinCastInfo from '@/components/library/LinCastInfo.vue';
import { ref } from 'vue';

export default {
  components: {
    Item,
    RecentPodcast,
    'lincast-info': LinCastInfo,
  },
  setup() {
    const recentPodcasts = ref([
      { title: 'Sample Podcast 1', key: 11 },
      { title: 'Sample Podcast 2', key: 12 },
      { title: 'Sample Podcast 3', key: 13 },
      { title: 'Sample Podcast 4', key: 14 },
      { title: 'Sample Podcast 5', key: 15 },
      { title: 'Sample Podcast 6', key: 16 },
    ]);

    return {
      recentPodcasts,
    };
  },
};
</script>

<style>
</style>
