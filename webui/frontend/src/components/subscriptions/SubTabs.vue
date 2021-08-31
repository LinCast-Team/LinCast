<template>
   <div class="flex justify-evenly border-b-2 border-gray-700">
      <button @click="activeTab = 'PodcastTab'" class="text-secondary-dt w-28 h-8" :class="{ 'activeTab':activeTab === 'PodcastTab' }">Podcasts</button>
      <button @click="activeTab = 'NewEpisodesTab'" class="text-secondary-dt w-28 h-8" :class="{ 'activeTab':activeTab === 'NewEpisodesTab' }">New Episodes</button>
    </div>
    <component :is="activeTab" />
</template>

<script lang='ts'>
import { ref } from 'vue';
import PodcastTab from './tabs/PodcastTab.vue';
import NewEpisodesTab from './tabs/NewEspisodesTab.vue';

export default {
  components: {
    PodcastTab,
    NewEpisodesTab,
  },
  setup() {
    const activeTab = ref('PodcastTab');

    return {
      activeTab,
    };
  },
};
</script>

<style lang='scss' scoped>
@import "@/assets/css/_palette.scss";

.activeTab {
  color: $primary-accent;
  border-color: $primary-accent;
  @apply border-b-2;
  @apply rounded-none;
}
</style>
