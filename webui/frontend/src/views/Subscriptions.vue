<template>
  <div class="font-sans flex flex-col">
    <lc-header
      title='Subscriptions'
      class="mt-6"
    />
    <button class="text-primary-dt font-semibold btn-bg mx-8 rounded py-0.5 my-2 ">
      <div class="flex flex-row justify-center">
        Listen All
        <div v-html="playIcon" class="my-auto mx-2"></div>
      </div>
    </button>
    <div class="overflow-y-auto">
      <!-- items here -->
      <!-- <list-item></list-item> -->
    </div>
  </div>
</template>

<script lang='ts'>
import LCHeader from '@/components/subscriptions/Header.vue';
// import ListItem from '@/components/subscriptions/ListItem.vue';
import { computed } from 'vue';
import feather from 'feather-icons';

export default {
  components: {
    'lc-header': LCHeader,
    // ListItem,
  },
  setup() {
    const playIcon = computed(() => feather.icons.play.toSvg({ 'stroke-width': 2.5, class: 'w-5 h-5' }));

    return {
      playIcon,
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
