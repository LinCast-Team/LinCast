<template>
  <div class="flex flex-col p-0.5">
      <img :src="imgSrc" :alt="title + '\'s artwork'" class="rounded-md w-28 h-28 mb-1">
      <h1 class="truncate text-primary-dt">{{ title }}</h1>
      <h3 class="truncate text-secondary-dt text-sm">{{ author }}</h3>
  </div>
</template>

<script lang='ts'>
import { defineComponent } from 'vue';
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
  },
});
</script>

<style>
</style>
