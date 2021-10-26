<template>
  <div class="font-sans flex flex-col overflow-y-auto">
    <div class="mt-6">
      <div class="flex flex-row justify-between my-2 mx-6">
        <h3 class="text-primary-dt text-xl font-semibold">Recently</h3>
        <div v-html="settingsIcon"></div>
      </div>
      <div v-if="recentPodcasts" class="flex flex-row my-2 overflow-x-auto">
        <recent-podcast
          v-for="p in recentPodcasts"
          :title="p.title"
          :key="p.ID"
          :imgSrc='p.imageURL'
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
import {
  ref,
  computed,
  defineComponent,
} from 'vue';
import feather from 'feather-icons';
import dayjs from 'dayjs';
import Item from '@/components/library/Item.vue';
import RecentPodcast from '@/components/library/RecentPodcast.vue';
import LinCastInfo from '@/components/library/LinCastInfo.vue';
import { SubscriptionsAPI } from '@/api';
import { Podcast } from '@/api/types';

export default defineComponent({
  components: {
    Item,
    RecentPodcast,
    'lincast-info': LinCastInfo,
  },
  setup() {
    const subsAPI = new SubscriptionsAPI();
    const recentPodcasts = ref(new Array<Podcast>());

    const setRecentPodcasts = async () => {
      const currentDate = dayjs();
      const previousDate = currentDate.subtract(30, 'day');

      const subscriptions = await subsAPI.getSubscriptions();
      const latestEps = await subsAPI.getLatestSubscriptionsEpisodes(previousDate.format('YYYY-MM-DD'), currentDate.format('YYYY-MM-DD'));
      const podcasts: Podcast[] = [];

      latestEps.forEach((ep) => {
        let parentPodcast = podcasts.find((p) => p.ID === ep.parentPodcastID);

        if (parentPodcast == null) {
          parentPodcast = subscriptions.find((p) => p.ID === ep.parentPodcastID);
          if (parentPodcast) {
            podcasts.push(parentPodcast);
          }
        }
      });

      recentPodcasts.value = podcasts;
    };

    setRecentPodcasts();

    const settingsIcon = computed(() => feather.icons.settings.toSvg({ 'stroke-width': 1.5, class: 'text-secondary-dt w-7 h-7 mx-4' }));

    return {
      recentPodcasts,
      settingsIcon,
    };
  },
});
</script>

<style>
</style>
