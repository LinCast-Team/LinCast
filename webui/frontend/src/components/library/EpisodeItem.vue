<template>
  <div class="flex flex-col w-full p-2">
      <div class="flex justify-between">
        <div class="flex">
          <img :src="imgSrc" :alt="title + '\'s artwork'" class=" rounded-md w-14">
          <div class=" pl-3">
            <h1 class="truncate text-secondary-dt">{{ title }}</h1>
            <h2 class="truncate text-secondary-dt">{{ author }}</h2>
          </div>
        </div>
        <h3 class="truncate text-secondary-dt mt-1">{{ update }}</h3>
        <p class="truncate text-secondary-dt mt-1">{{ duration }}</p>
      </div>
      <p class="truncate text-secondary-dt my-3">{{ resume }}</p>
      <div class="flex justify-between">
        <div class="flex">
          <span v-html="bookmarkIcon"></span>
          <span v-html="checkIcon"></span>
          <span v-html="arrowdowncircleIcon"></span>
          <span v-html="moreverticalIcon"></span>
        </div>
        <a class="flex align-middle py-0.5 px-3 border border-green-500 rounded-xl bg-green-600 text-white" href="#">Play <span class=" pl-1" v-html="playIcon"></span></a>
      </div>
  </div>
</template>

<script>
import defaultArtwork from '@/assets/resources/default_artwork.svg';
import feather from 'feather-icons';
import { computed } from 'vue';

export default {
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
    const bookmarkIcon = computed(() => feather.icons.bookmark.toSvg({ 'stroke-width': 1.5, class: 'text-secondary-dt w-7 h-7 ' }));
    const checkIcon = computed(() => feather.icons.check.toSvg({ 'stroke-width': 1.5, class: 'text-secondary-dt w-7 h-7 ' }));
    const playIcon = computed(() => feather.icons.play.toSvg({ 'stroke-width': 1.5, class: 'text-white w-6 ' }));
    const arrowdowncircleIcon = computed(() => feather.icons['arrow-down-circle'].toSvg({ 'stroke-width': 1.5, class: 'text-secondary-dt w-7 h-7 ' }));
    const moreverticalIcon = computed(() => feather.icons['more-vertical'].toSvg({ 'stroke-width': 1.5, class: 'text-secondary-dt w-7 h-7 ' }));

    return {
      bookmarkIcon,
      checkIcon,
      arrowdowncircleIcon,
      moreverticalIcon,
      playIcon,
    };
  },
};
</script>

<style>
</style>
