<template>
  <div id="nav" class="grid grid-cols-3 z-50 justify-items-center text-gray-50 bg-gradient-to-br from-gray-700 to-gray-900 px-4 md:px-16 lg:px-40">
    <router-link class="p-4 transition-colors duration-500 rounded-full" v-html="homeIcon" :to="{ name: 'Home' }"/>
    <router-link class="p-4 transition-colors duration-500 rounded-full" v-html="searchIcon" :to="{ name: 'Discover' }"/>
    <router-link class="p-4 transition-colors duration-500 rounded-full" v-html="bookmarkIcon" :to="{ name: 'Library' }"/>
  </div>
</template>

<script>
import { computed } from 'vue';
import feather from 'feather-icons';

export default {
  setup() {
    const homeIcon = computed(() => feather.icons.home.toSvg({ class: '' }));
    const searchIcon = computed(() => feather.icons.search.toSvg({ class: '' }));
    const bookmarkIcon = computed(() => feather.icons.bookmark.toSvg({ class: '' }));

    return {
      homeIcon,
      searchIcon,
      bookmarkIcon,
    };
  },
};
</script>

<style lang="scss" scoped>
.router-link-exact-active {
  color:#14B8A6;
}
</style>
