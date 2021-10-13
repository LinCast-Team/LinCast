<template>
  <div class="flex flex-col w-full py-3 px-4">
      <div class="flex justify-between">
        <div class="flex">
          <img :src="imgSrc" :alt="title + '\'s artwork'" class="rounded-md w-12 h-12 self-center">
          <div class="ml-3">
            <h1 class="truncate text-primary-dt text-xl">{{ title }}</h1>
            <h2 class="truncate text-secondary-dt">{{ author }}</h2>
          </div>
        </div>
        <p class="truncate text-secondary-dt mt-1 text-sm">{{ duration }}</p>
      </div>
      <p class="truncate text-primary-dt my-3">{{ resume }}</p>
      <div class="flex justify-between">
        <div class="flex space-x-2">
          <div v-html="bookmarkIcon"></div>
          <div v-html="checkIcon"></div>
          <div v-html="arrowdowncircleIcon"></div>
          <div v-html="moreverticalIcon"></div>
        </div>
        <a class="flex align-middle py-0.5 px-2.5 border rounded-2xl play-btn" href="#">Play<div class="ml-1 self-center" v-html="playIcon"></div></a>
      </div>
  </div>
</template>

<script lang='ts'>
import { defineComponent, computed } from 'vue';
import feather from 'feather-icons';
import defaultArtwork from '@/assets/resources/default_artwork.svg';

export default defineComponent({
  props: {
    imgSrc: {
      type: String,
      required: false,
      default: defaultArtwork,
    },
    title: {
      type: String,
      required: true,
    },
    author: {
      type: String,
      required: true,
    },
    resume: {
      type: String,
      required: true,
    },
    duration: {
      type: String,
      required: true,
    },
  },
  setup() {
    const bookmarkIcon = computed(() => feather.icons.bookmark.toSvg({ 'stroke-width': 1.2, class: 'text-secondary-dt w-7 h-7' }));
    const checkIcon = computed(() => feather.icons.check.toSvg({ 'stroke-width': 1.2, class: 'text-secondary-dt w-7 h-7' }));
    const playIcon = computed(() => feather.icons.play.toSvg({ 'stroke-width': 1.5, fill: 'currentColor', class: 'text-white w-5 h-5' }));
    const arrowdowncircleIcon = computed(() => feather.icons['arrow-down-circle'].toSvg({ 'stroke-width': 1.2, class: 'text-secondary-dt w-7 h-7' }));
    const moreverticalIcon = computed(() => feather.icons['more-vertical'].toSvg({ 'stroke-width': 1.2, class: 'text-secondary-dt w-7 h-7 ' }));

    return {
      bookmarkIcon,
      checkIcon,
      arrowdowncircleIcon,
      moreverticalIcon,
      playIcon,
    };
  },
});
</script>

<style lang='scss'>
@import "@/assets/css/_palette.scss";

.play-btn {
  background-color: $primary-accent;
  color: $text-primary;
  border-color: $primary-accent;

}
</style>
