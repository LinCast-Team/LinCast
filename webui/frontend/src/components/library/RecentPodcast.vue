<template>
  <div class="flex flex-col w-36 my-2 mx-4 p-1">
      <img :src="imgSrc" :alt="title + '\'s artwork'" class="h-32 w-32">
      <p class="truncate text-gray-400 mt-2">{{ title }}</p>
  </div>
</template>

<script>
import defaultArtwork from '@/assets/resources/default_artwork.svg';

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
  },
};
</script>

<style>
</style>
