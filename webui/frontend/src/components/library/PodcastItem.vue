<template>
  <li class="flex flex-col w-28 p-1 text-center float-left place-self-center">
      <img :src="imgSrc" :alt="title + '\'s artwork'" class="rounded-md w-24 h-24 m-1 self-center">
      <h1 class="truncate text-primary-dt">{{ title }}</h1>
      <h3 class="truncate text-secondary-dt text-sm">{{ author }}</h3>
  </li>
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
