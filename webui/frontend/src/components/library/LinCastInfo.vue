<template>
  <div class="flex flex-col items-center p-4 my-2 text-primary-dt">
    <img :src="lincastLogo" alt="LinCast logo" class="w-3/5" />
    <button class="rounded py-0.5 px-3 my-2 text-primary-dt btn-bg">
      Support Us
    </button>
    <div class="flex flex-row space-x-3 my-2">
      <a class="flex flex-row">
        <div v-html="githubIcon"></div>
        <p>GitHub</p>
      </a>
      <div class="flex flex-row">
        <div v-html="linkIcon"></div>
        <p>Site</p>
      </div>
      <div class="flex flex-row">
        <div v-html="twitterIcon"></div>
        <p>Twitter</p>
      </div>
    </div>
  </div>
</template>

<script>
import feather from 'feather-icons';
import { computed } from 'vue';
import lincastLogo from '@/assets/logo/logo_title.png';

export default {
  setup() {
    const githubIcon = computed(() => feather.icons.github.toSvg({ 'stroke-width': 1.5, class: 'w-6 h-6 mx-1' }));
    const linkIcon = computed(() => feather.icons.link.toSvg({ 'stroke-width': 1.5, class: 'w-6 h-6 mx-1' }));
    const twitterIcon = computed(() => feather.icons.twitter.toSvg({ 'stroke-width': 1.5, class: 'w-6 h-6 mx-1' }));

    return {
      lincastLogo,
      githubIcon,
      linkIcon,
      twitterIcon,
    };
  },
};
</script>

<style lang='scss' scoped>
@import "@/assets/css/_palette.scss";

.btn-bg {
  background-color: $primary-accent;
}
</style>
