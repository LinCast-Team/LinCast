<template>
  <div class="font-sans flex flex-col overflow-y-auto">
    <div class="mt-6">
      <div class="flex flex-row justify-between my-2 mx-6">
        <h3 class="text-primary-dt text-xl font-semibold">Recently</h3>
        <div v-html="settingsIcon"></div>
      </div>
      <div class="flex flex-row my-2 overflow-x-auto">
        <recent-podcast
          v-for="rp in recentPodcasts"
          :title="rp.title"
          :key="rp.key"
          imgSrc='https://picsum.photos/600'
        />
      </div>
      <div class="w-full bg-gray-500 opacity-25" style="height: 1px"></div>
      <div class="py-2 ">
        <item title="History" icon-class="book" />
        <item title="Listen Later" icon-class="clock" />
        <item title="Likes" icon-class="thumbs-up" />
      </div>
      <div class="w-full bg-gray-500 opacity-25" style="height: 1px"></div>
    </div>
    <lincast-info />
  </div>
</template>

<script lang='ts'>
import { ref, computed } from 'vue';
import feather from 'feather-icons';
import Item from '@/components/library/Item.vue';
import RecentPodcast from '@/components/library/RecentPodcast.vue';
import LinCastInfo from '@/components/library/LinCastInfo.vue';

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

    const settingsIcon = computed(() => feather.icons.settings.toSvg({ 'stroke-width': 1.5, class: 'text-secondary-dt w-7 h-7 mx-4' }));

    return {
      recentPodcasts,
      settingsIcon,
    };
  },
};
</script>

<style>
</style>
