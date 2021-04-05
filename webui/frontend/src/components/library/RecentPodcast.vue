<template>
  <div class="flex flex-col w-36 my-2 mx-4 p-1">
      <img :src="imgSrc" :alt="title + '\'s artwork'" class="h-32 w-32 rounded-md">
      <p class="truncate text-secondary-dt mt-1">{{ title }}</p>
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
