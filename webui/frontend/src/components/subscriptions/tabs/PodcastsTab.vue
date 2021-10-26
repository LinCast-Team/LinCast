<template>
  <div class="px-2 py-5">
    <ul v-if="subscriptions?.length > 0" class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 xl:grid-cols-10 2xl:grid-cols-12 gap-2 w-full">
      <podcast-item
        v-for="p in subscriptions"
        :key="p.id"
        :title="p.title"
        :author="p.authorName"
        :imgSrc="p.imageURL"
      />
    </ul>
    <div v-else class="mx-auto text-center text-secondary-dt">
      Looks like there is nothing here...
    </div>

    <!-- This should avoid content hidden by the navbar and the player -->
    <div class="h-sc" style="height: 16vh;"></div>
  </div>
</template>

<script lang='ts'>
import {
  defineComponent,
  inject,
  ref,
} from 'vue';
import PodcastItem from '@/components/library/PodcastItem.vue';
import { SubscriptionsAPI } from '@/api';
import { Podcast } from '@/api/types';

export default defineComponent({
  components: {
    'podcast-item': PodcastItem,
  },
  setup() {
    const subscriptions = inject('subscriptions', ref<Podcast[]>());
    const subsAPI = new SubscriptionsAPI();

    subsAPI.getSubscriptions()
      .then((res) => {
        subscriptions.value = res;
      })
      .catch((err) => {
        console.error(err);
      });

    return {
      subscriptions,
    };
  },
});
</script>

<style scoped>

</style>
