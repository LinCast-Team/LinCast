<template>
   <div class="flex justify-evenly sticky pt-3 top-0 border-b-2 border-gray-700 w-100 h-100">
      <button @click="onTabChange('PodcastsTab')" class="text-secondary-dt w-28 h-8" :class="{ 'activeTab':activeTab === 'PodcastsTab' }">Podcasts</button>
      <button @click="onTabChange('NewEpisodesTab')" class="text-secondary-dt w-28 h-8" :class="{ 'activeTab':activeTab === 'NewEpisodesTab' }">New Episodes</button>
    </div>
</template>

<script lang='ts'>
import {
  defineComponent,
  ref,
  SetupContext,
} from 'vue';

interface Data {
  [key: string]: unknown;
}

export default defineComponent({
  emits: [
    'tab-change',
  ],
  setup(_: Data, context: SetupContext): Data {
    const activeTab = ref('PodcastsTab');

    const onTabChange = (event: string) => {
      context.emit('tab-change', event);
      activeTab.value = event;
    };

    return {
      activeTab,
      onTabChange,
    };
  },
});
</script>

<style lang='scss' scoped>
div {
  background-color: $bg-primary;
}

.activeTab {
  color: $primary-accent;
  border-color: $primary-accent;
  @apply border-b-2;
  @apply rounded-none;
}
</style>
