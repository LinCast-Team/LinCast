<template>
  <div class="m-6 space-y-5">
    <for-you/>
    <subscriptions-preview/>
  </div>
</template>

<script lang="ts">
import SubscriptionsPreview from '@/components/home/SubscriptionsPreview.vue';
import ForYou from '@/components/home/ForYou.vue';

export default {
  components: {
    SubscriptionsPreview,
    ForYou,
  },
};
</script>
