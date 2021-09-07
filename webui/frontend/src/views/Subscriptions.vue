<template>
  <div class="font-sans flex flex-col">
    <lc-header
      title='Subscriptions'
      class="mt-6"
    />
    <SubTabs @tab-change="activeTab = $event"/>
    <component :is="activeTab" />
  </div>
</template>

<script lang='ts'>
import { ref } from 'vue';
import LCHeader from '@/components/subscriptions/Header.vue';
import SubTabs from '@/components/subscriptions/SubTabs.vue';
import PodcastsTab from '../components/subscriptions/tabs/PodcastsTab.vue';
import NewEpisodesTab from '../components/subscriptions/tabs/NewEspisodesTab.vue';

export default {
  components: {
    'lc-header': LCHeader,
    SubTabs,
    PodcastsTab,
    NewEpisodesTab,
  },
  setup() {
    const activeTab = ref('PodcastsTab');
    return {
      activeTab,
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
