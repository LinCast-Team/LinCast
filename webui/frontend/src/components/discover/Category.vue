<template>
  <div class="m-2 p-0.5 rounded-md bg-gradient-to-r from-green-700 to-green-600 w-64">
    <router-link tag="a" to="" class="p-3 w-64 text-gray-200 inline-block">{{ title }}</router-link>
  </div>
</template>

<script lang='ts'>
import { defineComponent } from 'vue';

export default defineComponent({
  props: {
    title: {
      type: String,
      required: true,
    },
  },
});
</script>

<style>
</style>
