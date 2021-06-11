<template>
  <div class="flex flex-col p-0.5">
      <img :src="imgSrc" :alt="title + '\'s artwork'" class=" rounded-md">
      <h1 class="truncate text-secondary-dt mt-1">{{ title }}</h1>
      <h3 class="truncate text-secondary-dt mt-1">{{ update }}</h3>
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
    update: {
      type: String,
      required: true,
    },
  },
};
</script>

<style>
</style>
