<template>
    <div class="m-2 p-2 flex">
        <img :src="artworkSrc" :alt="title + '\'s artwork'" class="w-12 h-12 md:w-20 md:h-20 rounded-xl mr-2">
        <div class="flex flex-col">
          <h3 class="text-gray-100 inline-block truncate font-medium text-lg">{{ title }}</h3>
          <h4 class="text-gray-200 inline-block truncate font-normal text-sm">{{ author }}</h4>
        </div>
    </div>
</template>
<script>
import defaultIcon from '@/assets/resources/default_artwork.svg';

export default {
  props: {
    title: {
      type: String,
      required: true,
    },
    author: {
      type: String,
      required: true,
    },
    artworkSrc: {
      type: String,
      required: false,
      default: defaultIcon,
    },
  },
  setup() {
    return {};
  },
};
</script>

<style>
</style>
