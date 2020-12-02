<template>
  <div id="nav">
    <router-link to="/">Home</router-link> |
    <router-link to="/about">About</router-link>
  </div>
  <router-view />
  <player
    :audioSrc="'https://wavesurfer-js.org/example/media/demo.wav'"
    :artworkSrc="'http://static-2.ivoox.com/canales/1/5/3/4/7691470744351_XXL.jpg'"
    :podcastTitle="'La tortulia podcast'"
    :episodeTitle="'India vs China'"
  />
</template>

<script>
import Player from '@/components/Player.vue';

export default {
  components: {
    Player,
  },
};
</script>

<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
</style>
