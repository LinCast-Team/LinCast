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
      :artworkSrc="'https://picsum.photos/1200'"
      :podcastTitle="'La tortulia podcast'"
      :episodeTitle="'India vs China'"
      :expanded="playerExpanded"
      :episodeDescription="'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam porttitor vitae velit ac rutrum. Etiam vitae ligula ac dui vestibulum dapibus. Sed fringilla nunc et volutpat euismod. Nullam suscipit, augue non mattis porttitor, magna mauris vehicula velit, ut tristique lacus arcu eu odio. Phasellus mauris nunc, ultricies sit amet leo at, suscipit sagittis metus. In condimentum nulla tristique, eleifend felis eget, dapibus tellus. Fusce tincidunt, turpis non euismod varius, nulla justo congue lectus, et vestibulum dolor purus tincidunt leo. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam iaculis vitae arcu sed rutrum. Donec elementum tempus cursus. Duis eu nisl pharetra, venenatis velit vitae, porttitor lectus. Nullam euismod imperdiet condimentum.'"
      @openRequest="openPlayer"
      @closeRequest="closePlayer"
    />
    <navigation-bar v-show="!playerExpanded" id="nav"/>
  </div>
</template>

<script lang='ts'>
import { ref } from 'vue';
import anime from 'animejs';
import Player from '@/components/Player.vue';
import NavigationBar from '@/components/NavigationBar.vue';

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
    };
  },
};
</script>

<style lang="scss">
@import "@/assets/css/_palette.scss";

body {
  @apply bg-primary-dt;
}

.bg-real-gray {
  background-color: #222529
}

// ========== Color Palette Dark Theme ===========

// Background color

.bg-primary-dt {
  background-color: $bg-primary;
}

.bg-secondary-dt {
  background-color: $bg-secondary;
}

.bg-accent-dt {
  background-color: $bg-accent;
}

// Text color

.text-primary-dt {
  color: $text-primary;
}

.text-secondary-dt {
  color: $text-secondary;
}

// Primary Accent

.primary-accent {
  color: $primary-accent;
}

.secondary-accent {
  color: $secondary-accent;
}

.link {
  color: $primary-accent;
}
</style>
