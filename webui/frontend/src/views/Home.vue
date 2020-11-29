<template>
  <div class="home">
    <div class="flex">
      <rounded-podcast-cover :podcastTitle='"Podcast 1"'/>
      <rounded-podcast-cover :podcastTitle='"Podcast 2"'/>
      <rounded-podcast-cover :podcastTitle='"Podcast 3"'/>
      <rounded-podcast-cover :podcastTitle='"Podcast 4"'/>
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import HelloWorld from '@/components/HelloWorld.vue'; // @ is an alias to /src
import RoundedPodcastCover from '@/components/RoundedPodcastCover.vue';

@Options({
  components: {
    HelloWorld,
    RoundedPodcastCover,
  },
})
export default class Home extends Vue {}
</script>
