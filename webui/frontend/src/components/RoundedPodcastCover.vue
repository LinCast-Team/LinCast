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
// import { ref } from "vue";

export default {
  props: {
    artworkSrc: {
      type: String,
      required: false,
      default:
        'http://static-2.ivoox.com/canales/1/5/3/4/7691470744351_XXL.jpg',
    },
    podcastTitle: {
      type: String,
      requied: true,
    },
  },
};
</script>>

<style>
</style>
