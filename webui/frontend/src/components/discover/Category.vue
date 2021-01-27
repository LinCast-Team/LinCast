<template>
  <div class="m-2 p-0.5 rounded-md bg-gradient-to-r from-indigo-700 to-green-700 w-64">
    <router-link tag="a" to="" class="p-3 w-64 text-gray-200 inline-block">{{ title }}</router-link>
  </div>
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      required: true,
    },
  },
  setup() {
    return {};
  },
};
</script>

<style>
</style>
