<template>
  <div class="flex p-2 py-5 justify-between">
    <div v-if="podcasts.length > 0">
      <podcast-item
        v-for="p in podcasts"
        :key="p.id"
        :title="p.title"
        :author="p.authorName"
        :imgSrc="p.imageURL"
      />
    </div>
    <div v-else class="mx-auto">
      <p class="self-center text-secondary-dt">Looks like there is nothing here...</p>
    </div>
  </div>
</template>

<script lang='ts'>
import { defineComponent, ref } from 'vue';
import PodcastItem from '@/components/library/PodcastItem.vue';
import { SubscriptionsAPI } from '@/api';
import { Podcast } from '@/api/types';

export default defineComponent({
  components: {
    'podcast-item': PodcastItem,
  },
  setup() {
    const podcasts = ref(new Array<Podcast>());
    const subsAPI = new SubscriptionsAPI();

    subsAPI.getSubscriptions()
      .then((res) => {
        podcasts.value = res;
      })
      .catch((err) => {
        console.error(err);
      });

    return {
      podcasts,
    };
  },
});
</script>

<style scoped>

</style>
