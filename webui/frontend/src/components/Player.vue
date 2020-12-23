<template>
<div
  class="flex z-50 shadow-lg font-sans transition-colors duration-500 text-center"
  :class="{
    'flex-row items-center bg-real-gray border-solid border-b border-black': !expanded,
    'flex-col w-full h-full tea-gradient overflow-y-auto': expanded,
  }"
> <!-- bg-gradient-to-br from-gray-700 to-gray-900 -->
  <h1 v-show="expanded" class="text-2xl my-8 text-gray-300">Playing Now</h1>

  <img
    id="player__podcast-artwork"
    :src="artworkSrc"
    :alt="podcastTitle + '\'s artwork'"
    class="self-center rounded-md shadow-lg"
    :class="{
      'w-60 h-60 md:w-64 md:h-64 mx-auto': expanded,
      'w-12 h-12 flex-none m-4': !expanded,
    }"
  >

  <div
    class=" text-center text-gray-100"
    :class="{
      'flex flex-row justify-around my-2 sm:mt-4 sm:mb-6 sm:mx-14 md:mx-20': expanded,
      'flex-grow cursor-pointer': !expanded,
    }"
    @click="if (!expanded) emitOpenEvent();"
  >
    <div v-show="expanded" v-html="share2Icon"></div>
    <div class="" :class="{ 'text-center justify-self-center w-3/5': expanded, 'text-left': !expanded }">
      <p
        id="player__podcast-title"
        class="truncate text-gray-100 uppercase"
        :class="{
          'font-extrabold text-xl': expanded,
          'font-bold text-base': !expanded,
        }"
      >{{ podcastTitle }}</p>
      <p
        id="player__episode-title"
        class="truncate text-gray-400"
        :class="{
          'text-sm': !expanded
        }"
      >{{ episodeTitle }}</p>
    </div>
    <div v-show="expanded" v-html="moreVerticalIcon"></div>
  </div>

  <audio ref="audioElement" :src="audioSrc" preload="auto"></audio>

  <div v-show="expanded" class="flex flex-col gap-2 justify-items-start mx-6 my-6">
    <div class="flex-grow bg-gradient-to-r from-gray-500 to-gray-700 rounded-md h-1 shadow-inner">
      <div class="rounded-md h-full w-0 shadow-inner" :style="'background-color: #14B8A6; width: ' + calculatedProgress  + '%;'"></div>
    </div>
    <div class="flex flex-row text-gray-400 font-bold text-sm bg-transparent mx-2 justify-between">
      <p>{{ currentTimeStr }}</p>
      <p>-{{ remainingTimeStr }}</p>
    </div>
  </div>

  <div
    id="player__buttons"
    class="bg-transparent transition-colors duration-500 text-gray-100"
    :class="{
      'grid grid-cols-5 items-center md:mx-20 py-4 px-1 sm:px-3 lg:px-1 xl:px-3 my-auto': expanded,
      'self-center': !expanded,
    }"
  >
    <button v-show="expanded" class="mx-auto rounded-full">
      <div v-html="skipBackIcon"></div>
    </button>
    <button v-show="expanded" @click="skipBackward(15)" class="mx-auto rounded-full">
      <div v-html="rotateCcwIcon"></div>
    </button>
    <button @click="playPause" class=" mx-4 rounded-full" :class="{ 'mx-auto': expanded, 'flex-none': !expanded }">
      <div v-if="expanded">
        <div v-if="!playing" v-html="playCirleIcon"></div>
        <div v-else v-html="pauseCirleIcon"></div>
      </div>
      <div v-else>
        <div v-if="!playing" v-html="playIcon"></div>
        <div v-else v-html="pauseIcon"></div>
      </div>
    </button>
    <button v-show="expanded" @click="skipForward(15)" class="mx-auto rounded-full">
      <div v-html="rotateCwIcon"></div>
    </button>
    <button v-show="expanded" class="mx-auto rounded-full">
      <div v-html="skipForwardIcon"></div>
    </button>
  </div>
  <div id="player__episode-description" v-if="episodeDescription !== ''" v-show="expanded" class="m-4 p-2 rounded-xl text-gray-200 bg-real-gray">
    {{ episodeDescription }}
  </div>
</div>
</template>

<script>
import {
  ref,
  computed,
  onMounted,
  onBeforeUnmount,
} from 'vue';
import feather from 'feather-icons';

// http://www.ivoox.com/tortulia-209-william-adams-parte-1_mf_60745571_feed_1.mp3

export default {
  props: {
    audioSrc: {
      type: String,
      required: true,
    },
    artworkSrc: {
      type: String,
      required: false,
      default: '',
    },
    podcastTitle: {
      type: String,
      required: true,
    },
    episodeTitle: {
      type: String,
      required: true,
    },
    episodeDescription: {
      type: String,
      required: false,
      default: '',
    },
    expanded: {
      type: Boolean,
      required: false,
      default: false,
    },
  },
  emits: ['open-request'],
  setup(_, context) {
    const playing = ref(false);
    const audioElement = ref(null);
    const currentTime = ref(0);
    const remainingTime = ref(0);
    const currentTimeStr = ref('00:00');
    const remainingTimeStr = ref('00:00');
    const duration = ref(0);

    const playCirleIcon = computed(() => feather.icons['play-circle'].toSvg({ 'stroke-width': 0.8, class: 'w-16 h-16 md:w-20 md:h-20' }));
    const pauseCirleIcon = computed(() => feather.icons['pause-circle'].toSvg({ 'stroke-width': 0.8, class: 'w-16 h-16 md:w-20 md:h-20' }));
    const playIcon = computed(() => feather.icons['play'].toSvg({ 'stroke-width': 1.0, class: 'w-9 h-9' })); /* eslint-disable-line */
    const pauseIcon = computed(() => feather.icons['pause'].toSvg({ 'stroke-width': 1.0, class: 'w-9 h-9' })); /* eslint-disable-line */
    const rotateCwIcon = computed(() => feather.icons['rotate-cw'].toSvg({ 'stroke-width': 1.5, class: 'w-8 h-8 md:w-12 md:h-12' }));
    const rotateCcwIcon = computed(() => feather.icons['rotate-ccw'].toSvg({ 'stroke-width': 1.5, class: 'w-8 h-8 md:w-12 md:h-12' }));
    const skipBackIcon = computed(() => feather.icons['skip-back'].toSvg({ 'stroke-width': 1.5, class: 'w-6 h-6 md:w-10 md:h-10' }));
    const skipForwardIcon = computed(() => feather.icons['skip-forward'].toSvg({ 'stroke-width': 1.5, class: 'w-6 h-6 md:w-10 md:h-10' }));
    const share2Icon = computed(() => feather.icons['share-2'].toSvg({ class: 'flex-none justify-self-center self-center mx-6 md:mx-14' }));
    const moreVerticalIcon = computed(() => feather.icons['more-vertical'].toSvg({ class: 'flex-none justify-self-center self-center mx-6 md:mx-12' }));

    const secsToMMSS = (secs) => {
      let minutes = Math.floor(secs / 60);
      let seconds = Math.floor(secs - (minutes * 60));

      if (minutes < 10) {
        minutes = `0${minutes}`;
      }
      if (seconds < 10) {
        seconds = `0${seconds}`;
      }

      return `${minutes}:${seconds}`;
    };

    const calculatedProgress = computed(() => (currentTime.value * 100) / duration.value);

    const playPause = () => {
      if (audioElement.value == null) {
        console.log('AudioElement null');
        return;
      }

      if (!playing.value) {
        console.log('Play clicked');
        audioElement.value.play();
      } else {
        console.log('Pause clicked');
        audioElement.value.pause();
      }
    };

    const skipBackward = (secs) => {
      if (audioElement.value == null) {
        return;
      }

      if (currentTime.value <= secs) {
        audioElement.value.currentTime = 0;
        currentTime.value = 0;
      } else {
        audioElement.value.currentTime -= secs;
        currentTime.value -= secs;
      }
    };

    const skipForward = (secs) => {
      if (audioElement.value == null) {
        return;
      }

      if ((currentTime.value + secs) >= duration.value) {
        audioElement.value.currentTime = duration.value;
        currentTime.value = duration.value;
      } else {
        audioElement.value.currentTime += secs;
        currentTime.value += secs;
      }
    };

    const updateRemaining = () => {
      remainingTime.value = Math.floor(duration.value) - Math.floor(currentTime.value);
      remainingTimeStr.value = secsToMMSS(remainingTime.value);
    };

    const updateDuration = () => {
      duration.value = audioElement.value.duration;
      updateRemaining();
    };

    const updateCurrentAndRemaining = () => {
      currentTime.value = audioElement.value.currentTime;
      currentTimeStr.value = secsToMMSS(currentTime.value);
      updateRemaining();
    };

    const setPlaying = () => { playing.value = true; };

    const setPaused = () => { playing.value = false; };

    const emitOpenEvent = () => {
      context.emit('open-request');
    };

    onMounted(() => {
      audioElement.value.addEventListener('durationchange', updateDuration);
      audioElement.value.addEventListener('timeupdate', updateCurrentAndRemaining);
      audioElement.value.addEventListener('play', setPlaying);
      audioElement.value.addEventListener('pause', setPaused);
    });

    onBeforeUnmount(() => {
      audioElement.value.removeEventListener('durationchange', updateDuration);
      audioElement.value.removeEventListener('timeupdate', updateCurrentAndRemaining);
      audioElement.value.removeEventListener('play', setPlaying);
      audioElement.value.removeEventListener('pause', setPaused);
    });

    return {
      // Icons
      playCirleIcon,
      pauseCirleIcon,
      playIcon,
      pauseIcon,
      rotateCwIcon,
      rotateCcwIcon,
      skipBackIcon,
      skipForwardIcon,
      share2Icon,
      moreVerticalIcon,

      // Player functionality
      audioElement,
      playing,
      currentTime,
      currentTimeStr,
      remainingTimeStr,
      calculatedProgress,
      playPause,
      skipForward,
      skipBackward,
      emitOpenEvent,
    };
  },
};
</script>

<style lang="scss" scoped>
.tea-gradient {
  background-image: linear-gradient(to bottom right, #004D40, #000, #000, #000);
}
</style>
