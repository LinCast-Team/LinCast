<template>
  <!-- <div id="nav">
    <router-link to="/">Home</router-link> |
    <router-link to="/about">About</router-link>
  </div> -->
  <router-view />
  <div id="player-container" class="fixed bottom-0 right-0 left-0 flex flex-col">
    <player
      id="player"
      :expanded="playerExpanded"
      @openRequest="openPlayer"
      @closeRequest="closePlayer"
    />
    <navigation-bar v-show="!playerExpanded" id="nav"/>
  </div>
</template>

<script lang='ts'>
import { ref, provide } from 'vue';
import anime from 'animejs';
import Player from '@/components/Player.vue';
import NavigationBar from '@/components/NavigationBar.vue';
import { Podcast } from './api/types';
import { SubscriptionsAPI } from './api';

export default {
  components: {
    Player,
    NavigationBar,
  },
  setup() {
    const subscriptions = ref<Podcast[]>();
    const subsAPI = new SubscriptionsAPI();
    provide('subscriptions', subscriptions);

    subsAPI.getSubscriptions()
      .then((subs) => {
        subscriptions.value = subs;
      })
      .catch((err) => {
        console.error(err);
      });

    const playerExpanded = ref(false);

    const openPlayer = () => {
      console.log('Player open request launched');

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
          // top: 0,
          translateY: 0,
          begin: () => {
            playerExpanded.value = true;
            const el = document.getElementById('player-container');
            if (el != null) {
              el.style.top = '0px';
            }
          },
        }, '-=300');
    };

    const closePlayer = () => {
      const tl = anime.timeline({
        targets: '#player-container',
        easing: 'easeOutExpo',
        duration: 700,
      });

      tl
        .add({
          opacity: [1, 0],
          translateY: 500,
        })
        .add({
          opacity: [0, 1],
          translateY: 0,
          begin: () => {
            playerExpanded.value = false;
            const el = document.getElementById('player-container');
            if (el != null) {
              el.style.top = 'auto';
            }
          },
        }, '-=300');
    };

    return {
      playerExpanded,
      openPlayer,
      closePlayer,
      subscriptions,
    };
  },
};
</script>

<style lang="scss">
body {
  @apply bg-primary-dt;
}

.bg-real-gray {
  background-color: #222529
}


</style>
