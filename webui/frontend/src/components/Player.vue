<template>
<div
  class="flex z-50 shadow-lg font-sans transition-colors duration-500 text-center"
  :class="{
    'flex-row items-center bg-real-gray border-solid border-b border-black': !expanded,
    'flex-col w-full h-full tea-gradient overflow-y-auto': expanded,
  }"
> <!-- bg-gradient-to-br from-gray-700 to-gray-900 -->
  <div v-show="expanded" @click="emitCloseEvent" class="flex-none rounded-lg bg-gray-500 w-1/5 m-auto my-5 h-1 shadow-md"></div>
  <h1 v-show="expanded" class="text-2xl mb-8 mt-2 text-gray-300">Playing Now</h1>

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
      'w-3/5 cursor-pointer': !expanded,
    }"
    @click="if (!expanded) emitOpenEvent();"
  >
    <div v-show="expanded" v-html="share2Icon" class="self-center"></div>
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
    <div v-show="expanded" v-html="moreVerticalIcon" class="self-center"></div>
  </div>

  <div v-show="expanded" class="flex flex-col gap-2 justify-items-start mx-6 my-6">
    <div class="flex-grow bg-gradient-to-r from-gray-500 to-gray-700 rounded-md h-1 shadow-inner flex">
      <div class="rounded-md h-full w-0 shadow-inner" :style="'background-color: #14B8A6; width: ' + calculatedProgress  + '%;'"></div>
      <div class="h-4 w-4 rounded-full border border-black relative -left-2" style="background-color: #14B8A6; top: -6px;"></div>
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
        <div v-if="!playing">
          <svg v-if="!loading" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="0.8" stroke-linecap="round" stroke-linejoin="round" class="w-16 h-16 md:w-20 md:h-20"><circle cx="12" cy="12" r="10"></circle><polygon points="10 8 16 12 10 16 10 8" fill="black" stroke="black"></polygon></svg>
          <svg v-else class="animate-spin w-16 h-16 md:w-20 md:h-20" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>
        <div v-else>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="0.8" stroke-linecap="round" stroke-linejoin="round" class="w-16 h-16 md:w-20 md:h-20"><circle cx="12" cy="12" r="10"></circle><line x1="10" y1="15" x2="10" y2="9" stroke="black" stroke-width="2"></line><line x1="14" y1="15" x2="14" y2="9" stroke="black" stroke-width="2"></line></svg>
        </div>
      </div>
      <div v-else>
        <div v-if="!playing">
          <svg v-if="!loading" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-9 h-9"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>
          <svg v-else class="animate-spin w-9 h-9" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>
        <div v-else>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round" class="w-9 h-9"><rect x="6" y="4" width="4" height="16" stroke="currentColor" fill="currentColor"></rect><rect x="14" y="4" width="4" height="16" stroke="currentColor" fill="currentColor"></rect></svg>
        </div>
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
import { Howl } from 'howler';

// http://www.ivoox.com/tortulia-209-william-adams-parte-1_mf_60745571_feed_1.mp3

export default {
  props: {
    audioSrc: {
      type: Array,
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
    autoplay: {
      type: Boolean,
      required: false,
      default: false,
    },
  },
  emits: ['open-request', 'close-request'],
  setup(props, context) {
    const playing = ref(false);
    const currentTime = ref(0);
    const remainingTime = ref(0);
    const currentTimeStr = ref('00:00');
    const remainingTimeStr = ref('00:00');
    const duration = ref(0);
    const volume = ref(1.0);
    const intervalID = ref(0);
    const loading = ref(true);

    const audio = ref(new Howl({
      src: props.audioSrc,
      autoplay: props.autoplay,
      volume: volume.value,
      loop: false,
      html5: true,
    }));

    // const currentTime = computed({
    //   get: () => audio.value.seek(),
    //   set: (val) => audio.value.seek(val),
    // });

    // const duration = computed(() => audio.value.duration());

    const rotateCwIcon = computed(() => feather.icons['rotate-cw'].toSvg({ 'stroke-width': 1.5, class: 'w-8 h-8 md:w-12 md:h-12' }));
    const rotateCcwIcon = computed(() => feather.icons['rotate-ccw'].toSvg({ 'stroke-width': 1.5, class: 'w-8 h-8 md:w-12 md:h-12' }));
    const skipBackIcon = computed(() => feather.icons['skip-back'].toSvg({ 'stroke-width': 1.5, class: 'w-6 h-6 md:w-10 md:h-10' }));
    const skipForwardIcon = computed(() => feather.icons['skip-forward'].toSvg({ 'stroke-width': 1.5, class: 'w-6 h-6 md:w-10 md:h-10' }));
    const share2Icon = computed(() => feather.icons['share-2'].toSvg({ class: 'mx-6 md:mx-14' }));
    const moreVerticalIcon = computed(() => feather.icons['more-vertical'].toSvg({ class: 'mx-6 md:mx-12' }));

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

    // const currentTimeStr = computed(() => secsToMMSS(currentTime.value));

    const calculatedProgress = computed(() => (currentTime.value * 100) / duration.value);

    const playPause = () => {
      if (!audio.value.state() === 'loaded') {
        return;
      }

      if (!playing.value) {
        audio.value.play();
      } else {
        audio.value.pause();
      }
    };

    const updateRemaining = () => {
      remainingTime.value = Math.floor(duration.value) - Math.floor(currentTime.value);
      remainingTimeStr.value = secsToMMSS(remainingTime.value);
    };

    const updateDuration = () => {
      duration.value = audio.value.duration();
      updateRemaining();
    };

    const updateCurrentAndRemaining = () => {
      currentTime.value = audio.value.seek();
      currentTimeStr.value = secsToMMSS(currentTime.value);
      updateRemaining();
    };

    const skipBackward = (secs) => {
      if (currentTime.value <= secs) {
        // audioElement.value.currentTime = 0;
        audio.value.seek(0);
      } else {
        // audioElement.value.currentTime -= secs;
        audio.value.seek(currentTime.value - secs);
      }
      updateCurrentAndRemaining();
    };

    const skipForward = (secs) => {
      if ((currentTime.value + secs) >= duration.value) {
        // audioElement.value.currentTime = duration.value;
        audio.value.seek(duration.value);
      } else {
        // audioElement.value.currentTime += secs;
        audio.value.seek(currentTime.value + secs);
      }
      updateCurrentAndRemaining();
    };

    const setPlaying = () => { playing.value = true; };

    const setPaused = () => { playing.value = false; };

    const emitOpenEvent = () => {
      context.emit('open-request');
    };

    const emitCloseEvent = () => {
      context.emit('close-request');
    };

    onMounted(() => {
      audio.value.on('load', () => {
        updateDuration();
        loading.value = false;
      });

      audio.value.on('loaderror', () => {
        // todo
      });

      audio.value.on('play', () => {
        setPlaying();
        updateCurrentAndRemaining();
        intervalID.value = setInterval(() => {
          updateCurrentAndRemaining();
        }, 500);
      });

      audio.value.on('playerror', () => {
        setPaused();
        clearInterval(intervalID.value);
      });

      audio.value.on('pause', () => {
        setPaused();
        clearInterval(intervalID.value);
      });

      audio.value.on('end', () => {
        setPaused();
        clearInterval(intervalID.value);
      });

      audio.value.on('stop', () => {
        setPaused();
        clearInterval(intervalID.value);
      });

      audio.value.on('seek', () => {
        // todo
      });

      audio.value.on('unlock', () => {
        // todo
      });
    });

    onBeforeUnmount(() => {
      audio.value.off();
    });

    return {
      // Icons
      rotateCwIcon,
      rotateCcwIcon,
      skipBackIcon,
      skipForwardIcon,
      share2Icon,
      moreVerticalIcon,

      // Player functionality
      audio,
      playing,
      loading,
      currentTime,
      currentTimeStr,
      remainingTimeStr,
      calculatedProgress,
      playPause,
      skipForward,
      skipBackward,
      emitOpenEvent,
      emitCloseEvent,
    };
  },
};
</script>

<style lang="scss" scoped>
.tea-gradient {
  background-image: linear-gradient(to bottom right, #004D40, #000, #000, #000);
}
</style>
