<template>
  <!-- <div id="nav">
    <router-link to="/">Home</router-link> |
    <router-link to="/about">About</router-link>
  </div> -->
  <router-view />
  <div id="player-container" class="fixed bottom-0 right-0 left-0 flex flex-col">
    <player
      id="player"
      :audioSrc="'http://www.ivoox.com/tortulia-209-william-adams-parte-1_mf_60745571_feed_1.mp3'"
      :artworkSrc="'http://static-2.ivoox.com/canales/1/5/3/4/7691470744351_XXL.jpg'"
      :podcastTitle="'La tortulia podcast'"
      :episodeTitle="'India vs China'"
      :expanded="playerExpanded"
      :episodeDescription="'Hello world'"
      @openRequest="openPlayer"
    />
    <navigation-bar v-show="!playerExpanded" id="nav"/>
  </div>
</template>

<script>
import { ref } from 'vue';
import Player from '@/components/Player.vue';
import NavigationBar from '@/components/NavigationBar.vue';
import anime from 'animejs/lib/anime.es';

export default {
  components: {
    Player,
    NavigationBar,
  },
  setup() {
    const playerExpanded = ref(false);

    const openPlayer = () => {
      const tl = anime.timeline({
        targets: '#player-container',
        easing: 'easeOutExpo',
        duration: 700,
      });

      tl
        .add({
          opacity: [1, 0],
          translateY: 100,
        })
        .add({
          opacity: [0, 1],
          top: 0,
          translateY: 0,
          // easing: 'spring(1, 80, 10, 0)',
          begin: () => {
            playerExpanded.value = true;
          },
        }, '-=300');
    };

    return { playerExpanded, openPlayer };
  },
};
</script>

<style lang="scss">
body {
  @apply bg-black;
}

.bg-real-gray {
  background-color: #222529;
}

// .hide-enter-active,
// .hide-leave-active {
//   transition: all 1s ease;
// }

// .hide-enter-from,
// .hide-leave-to {
//   transform: translateY(100px);
// }
</style>
