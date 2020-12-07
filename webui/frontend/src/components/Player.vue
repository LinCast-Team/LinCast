<template>
<div
  class="flex z-50 shadow-lg font-sans border-solid border-b-2 transition-colors duration-500 border-gray-200 bg-gray-50 dark:border-gray-700 dark:bg-gray-900"
  :class="{
    'rounded-t-md flex-row items-center border-t-2': !expanded,
    'flex-col w-full h-full': expanded,
  }"
>
  <h1 v-show="expanded" class="text-2xl my-8 dark:text-gray-300">Playing Now</h1>

  <img
    :src="artworkSrc"
    :alt="podcastTitle + '\'s artwork'"
    class="self-center rounded-md shadow-lg"
    :class="{
      'w-60 h-60 md:w-64 md:h-64 mx-auto': expanded,
      'w-12 h-12 flex-none m-4': !expanded,
    }"
  >

  <div
    class=" text-center dark:text-gray-100"
    :class="{
      'flex flex-row justify-around mt-4 mb-6 sm:mx-14 md:mx-20': expanded,
      'flex-grow': !expanded,
    }"
  >
    <div v-show="expanded" v-html="share2Icon"></div>
    <div class="" :class="{ 'flex-grow text-center justify-self-center': expanded, 'text-left': !expanded }">
      <p
        class="truncate text-black dark:text-gray-100 uppercase "
        :class="{
          'font-extrabold text-xl': expanded,
          'font-bold text-base': !expanded,
        }"
      >{{ podcastTitle }}</p>
      <p
        class="truncate text-gray-500 dark:text-gray-400"
        :class="{
          'text-sm': !expanded
        }"
      >{{ episodeTitle }}</p>
    </div>
    <div v-show="expanded" v-html="moreVerticalIcon"></div>
  </div>

  <div v-show="expanded" id="waveform" class="bg-transparent mx-4"></div>

  <div
    class="bg-transparent text-black transition-colors duration-500 dark:text-gray-100    "
    :class="{
      'grid grid-cols-5 items-center md:mx-20 py-4 px-1 sm:px-3 lg:px-1 xl:px-3 my-auto': expanded,
      'self-center': !expanded,
    }"
  >
    <button v-show="expanded" class="mx-auto rounded-full">
      <div v-html="skipBackIcon"></div>
    </button>
    <button v-show="expanded" @click="wavesurfer.skipBackward(15)" class="mx-auto rounded-full">
      <div v-html="rotateCcwIcon"></div>
    </button>
    <button @click="wavesurfer.playPause()" class=" mx-4 rounded-full" :class="{ 'mx-auto': expanded, 'flex-none': !expanded }">
      <!-- <span v-if="expanded" class="w-16 h-16 md:w-20 md:h-20" stroke-width="0.8" data-feather="play-circle"></span> -->
      <div v-if="expanded">
        <div v-if="!playing" v-html="playCirleIcon"></div>
        <div v-else v-html="pauseCirleIcon"></div>
      </div>
      <div v-else>
        <div v-if="!playing" v-html="playIcon"></div>
        <div v-else v-html="pauseIcon"></div>
      </div>
    </button>
    <button v-show="expanded" @click="wavesurfer.skipForward(15)" class="mx-auto rounded-full">
      <div v-html="rotateCwIcon"></div>
    </button>
    <button v-show="expanded" class="mx-auto rounded-full">
      <div v-html="skipForwardIcon"></div>
    </button>
  </div>
</div>
</template>

<script>
import { onMounted, ref, computed } from 'vue';
import WaveSurfer from 'wavesurfer.js';
import feather from 'feather-icons';
// import axios from 'axios';

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
  setup(props) {
    const wavesurfer = ref();
    const playing = ref(false);

    onMounted(() => {
      wavesurfer.value = WaveSurfer.create({
        container: '#waveform',
        waveColor: '#99F6E4',
        progressColor: '#14B8A6',
        barWidth: 3,
        barRadius: 3,
        barMinHeight: 1,
        cursorWidth: 0,
        height: 100,
        barGap: 3,
        responsive: true,
      });

      wavesurfer.value.load(props.audioSrc);

      let updateIntervalID;
      let previousProgress;

      wavesurfer.value.on('play', () => {
        // TODO Call the API to update progress using `wavesurfer.getCurrentTime()`.
        // If the current time remains the same, the request should not be sent.
        console.log('Playing');
        playing.value = true;

        updateIntervalID = setInterval(() => {
          const progress = wavesurfer.value.getCurrentTime();
          if (previousProgress !== progress) {
            console.log('Updating progress:', progress);
            previousProgress = progress;
          }
        }, 1000);
      });

      wavesurfer.value.on('pause', () => {
        // TODO Stop calling the API to update progress.
        console.log('Paused');
        playing.value = false;
        clearInterval(updateIntervalID);
      });

      wavesurfer.value.on('seek', (newPosition) => {
        console.log('New position on player\'s cursor:', newPosition);
        // TODO Call the API to update progress using `wavesurfer.getCurrentTime()`.
        // If the current time remains the same, the request should not be sent.
      });

      wavesurfer.value.on('loading', (progress) => {
        console.log(`Loading audio: ${progress}%`);
        // TODO Show the progress to the user.
      });

      wavesurfer.value.on('finish', () => {
        // TODO Load the next episode and play it.
        console.log('Audio completely played');
      });

      wavesurfer.value.on('destroy', () => {
        // TODO Stop calling the API to update progress.
        console.log('Wavesurfer instance destroyed');
        playing.value = false;
        clearInterval(updateIntervalID);
      });

      wavesurfer.value.on('error', (err) => {
        console.error(err);
      });
    });

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

    return {
      wavesurfer,
      playing,
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
    };
  },
};
</script>

<style>
</style>
