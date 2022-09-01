<template>
<div class="px-2 py-5">
  <div v-if="episodes?.length > 0">
    <h2 class="p-2 pt-3 text-secondary-dt">Today</h2>
      <!-- FIX Some props do not match with the fields of the object -->
    <ol class="flex flex-col flex-grow">
      <episode-item
        v-for="e in episodes"
        :title="e.title"
        :imgSrc="getPodcastArtwork(e.id)"
        :author="getPodcastName(e.id)"
        :resume="e.description"
        :duration="e.duration"
        :key="e.id"
      />
    </ol>
  </div>
  <div v-else class="mx-auto">
    <p class="self-center text-secondary-dt">Looks like there is nothing here...</p>
  </div>
</div>
</template>

<script lang='ts'>
import { defineComponent } from 'vue';
import { storeToRefs } from 'pinia';
import EpisodeItem from '@/components/library/EpisodeItem.vue';
import { useSubscriptionsStore } from '@/store/subscriptions';

export default defineComponent({
  components: {
    EpisodeItem,
  },
  setup() {
    const store = useSubscriptionsStore();
    const { podcasts, episodes } = storeToRefs(store);

    const getPodcastName = (id: number): string | undefined => {
      const podcast = podcasts.value.find((p) => p.ID === id);
      return podcast?.title;
    };

    const getPodcastArtwork = (id: number): string | undefined => {
      const podcast = podcasts.value.find((p) => p.ID === id);
      return podcast?.imageURL;
    };

    return {
      episodes,
      podcasts,
      getPodcastName,
      getPodcastArtwork,
    };
  },
});
</script>

<style scoped>

</style>
