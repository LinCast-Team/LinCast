<template>
  <!-- <div id="nav">
    <router-link to="/">Home</router-link> |
    <router-link to="/about">About</router-link>
  </div> -->
  <!-- <button @click="playerExpanded = !playerExpanded" class="fixed top-0 z-50">Expandir/Contraer</button> -->
  <router-view />
  <div class="fixed bottom-0 right-0 left-0 flex flex-col" :class="{ 'top-0': playerExpanded }">
    <player
      :audioSrc="'http://www.ivoox.com/tortulia-209-william-adams-parte-1_mf_60745571_feed_1.mp3'"
      :artworkSrc="'http://static-2.ivoox.com/canales/1/5/3/4/7691470744351_XXL.jpg'"
      :podcastTitle="'La tortulia podcast'"
      :episodeTitle="'India vs China'"
      :expanded="playerExpanded"
    />
    <navigation-bar/>
  </div>
</template>

<script>
import { ref } from 'vue';
import Player from '@/components/Player.vue';
import NavigationBar from '@/components/NavigationBar.vue';

export default {
  components: {
    Player,
    NavigationBar,
  },
  setup() {
    const playerExpanded = ref(true);

    return { playerExpanded };
  },
};
</script>

<style lang="scss">
</style>
