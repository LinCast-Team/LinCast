<template>
<div class="flex flex-col border-2 border-solid border-gray-100 bg-gray-50 rounded-t-xl z-100 shadow-lg font-sans" :class="{ 'fixed bottom-0 right-0 left-0 top-0': !expanded }">
  <h1 class="text-2xl my-8">Playing Now</h1>
  <div class="flex flex-col text-center items-center mb-6">
    <img :src="artworkSrc" :alt="podcastTitle + '\'s artwork'" class="w-64 h-64 md:w-80 md:h-80 mx-1 rounded-md shadow-lg">
    <div class="grid grid-cols-4 grid-flow-col gap-2 items-center justify-items-center my-4 md:my-8">
      <div class="" data-feather="share-2"></div>
      <div class="col-span-2 text-center">
        <p class="truncate text-black uppercase font-extrabold text-xl">{{ podcastTitle }}</p>
        <p class="truncate text-gray-500">{{ episodeTitle }}</p>
      </div>
      <div class="" data-feather="more-vertical"></div>
    </div>
  </div>
  <div id="waveform" class="bg-gray-50 mx-4"></div>
  <div class="bg-gray-50 text-black transition-colors duration-500 dark:bg-gray-900 dark:text-white py-4 px-1 sm:px-3 lg:px-1 xl:px-3 my-auto md:mx-20 grid grid-cols-5 items-center">
    <button class="mx-auto rounded-full">
      <span class="w-6 h-6 md:w-10 md:h-10" data-feather="skip-back" stroke-width="1.5"></span>
    </button>
    <button @click="wavesurfer.skipBackward(15)" class="mx-auto rounded-full">
      <span class="w-8 h-8 md:w-12 md:h-12" data-feather="rotate-ccw" stroke-width="1.5"></span>
    </button>
    <button @click="wavesurfer.playPause()" class="mx-auto rounded-full">
      <span class="w-16 h-16 md:w-20 md:h-20" stroke-width="0.8" data-feather="play-circle"></span>
    </button>
    <button @click="wavesurfer.skipForward(15)" class="mx-auto rounded-full">
      <span class="w-8 h-8 md:w-12 md:h-12" data-feather="rotate-cw" stroke-width="1.5"></span>
    </button>
    <button class="mx-auto rounded-full">
      <span class="w-6 h-6 md:w-10 md:h-10" data-feather="skip-forward" stroke-width="1.5"></span>
    </button>
  </div>
</div>
</template>

<script>
import { onMounted, ref } from 'vue';
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

    onMounted(() => {
      feather.replace();

      wavesurfer.value = WaveSurfer.create({
        container: '#waveform',
        waveColor: '#99F6E4',
        progressColor: '#14B8A6',
        barWidth: 3,
        barRadius: 3,
        cursorWidth: 0,
        height: 100,
        barGap: 3,
      });

      wavesurfer.value.load(props.audioSrc);

      wavesurfer.value.on('play', () => {
        // TODO Call the API to update progress using `wavesurfer.getCurrentTime()`.
        // If the current time remains the same, the request should not be sent.
        console.log('Playing');
      });

      wavesurfer.value.on('pause', () => {
        // TODO Stop calling the API to update progress.
        console.log('Paused');
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
      });

      wavesurfer.value.on('error', (err) => {
        console.error(err);
      });
    });

    return { wavesurfer };
  },
};
</script>

<style>
</style>
