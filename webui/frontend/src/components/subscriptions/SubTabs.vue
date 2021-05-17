<template>
   <div class="flex justify-evenly border-b-2 border-gray-700">
      <button @click="activeTab = 'PodcastTab'" class="primary-accent w-28 h-8 border-b-2 border-green-600 rounded-none" v-bind:class="['tab-button', { active: currentTab === tab }]">Podcasts</button>
      <button @click="activeTab = 'NewEpisodesTab'" class="text-secondary-dt w-28 h-8">New Episodes</button>
    </div>
    <component :is="activeTab" />
</template>

<script>
import PodcastTab from './tabs/PodcastTab.vue';
import NewEpisodesTab from './tabs/NewEspisodesTab.vue';

export default {
  components: {
    PodcastTab,
    NewEpisodesTab,
  },
  data() {
    return {
      activeTab: 'PodcastTab',
    };
  },
};
</script>

<style scoped>

.activeTab {
  color:#4DB6AC;
  border-color: #4DB6AC;
  @apply border-b-2;
  @apply border-green-600;
  @apply rounded-none;
}

</style>
