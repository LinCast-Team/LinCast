<template>
  <div class="flex flex-col items-center justify-center w-24 h-40 m-1">
    <img
      :src="artworkSrc"
      :alt="podcastTitle + '\'s artwork'"
      class="rounded-full w-24 h-24"
    />
    <p class="text-gray-600 font-sans">{{ podcastTitle }}</p>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import defaultIcon from '@/assets/resources/default_artwork.svg';

export default defineComponent({
  props: {
    artworkSrc: {
      type: String,
      required: false,
      default: defaultIcon,
    },
    podcastTitle: {
      type: String,
      requied: true,
    },
  },
});
</script>>

<style>
</style>
