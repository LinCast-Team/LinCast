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
      :episodeDescription="'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed ut mauris id felis ultricies dignissim a id arcu. Nam quis orci sed lectus pretium mollis ac quis quam. Quisque ultricies tempus auctor. Cras id urna varius, tincidunt magna non, venenatis magna. Pellentesque non aliquam mi, non congue ex. Nunc suscipit lacus non arcu euismod, et consectetur nunc vulputate. Aenean neque purus, scelerisque id placerat vel, sollicitudin ut enim. Donec mollis erat eget pretium gravida. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Nulla ullamcorper vitae arcu nec scelerisque. Aliquam gravida velit vitae elit mollis convallis. Nunc at odio nec justo semper tempus eget at quam. Proin eget nunc mattis, varius nunc vitae, accumsan turpis. Sed felis nunc, convallis at vehicula in, finibus eget turpis. Suspendisse id ante sem. '"
      @openRequest="playerExpanded = true"
    />
    <transition name="hide">
      <navigation-bar :hidden="playerExpanded"/>
    </transition>
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
    const playerExpanded = ref(false);

    return { playerExpanded };
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

.hide-enter-active,
.hide-leave-active {
  transition: all .5s ease;
}

.hide-enter-from,
.hide-leave-to {
  transform: translateY(50px);
}
</style>
