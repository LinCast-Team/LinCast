<template>
<div
  class="flex z-50 shadow-lg font-sans transition-colors duration-500 text-center text-secondary-dt"
  :class="{
    'flex-row items-center bg-primary-dt border-solid border-b border-gray-800': !expanded,
    'flex-col w-full h-full bg-primary-dt overflow-y-auto': expanded,
  }"
> <!-- bg-gradient-to-br from-gray-700 to-gray-900 -->
  <div v-show="expanded" @click="emitCloseEvent" class="flex-none rounded-lg bg-accent-dt opacity-40 w-1/5 m-auto my-5 h-1 shadow-md"></div>
  <h1 v-show="expanded" class="text-2xl mb-8 mt-2 text-primary-dt">Playing Now</h1>

  <img
    id="player__podcast-artwork"
    :src="artworkSrc || defaultArtwork"
    :alt="podcastTitle + '\'s artwork'"
    class="self-center rounded-md shadow-lg"
    :class="{
      'w-60 h-60 md:w-64 md:h-64 mx-auto': expanded,
      'w-12 h-12 flex-none m-4': !expanded,
      'opacity-40': !artworkSrc,
    }"
  >

  <div
    class=" text-center"
    :class="{
      'flex flex-row justify-around my-2 sm:mt-4 sm:mb-6 sm:mx-14 md:mx-20': expanded,
      'w-3/5 cursor-pointer': !expanded,
    }"
    @click="if (!expanded) emitOpenEvent();"
  >
    <div v-show="expanded" v-html="share2Icon" class="self-center text-primary-dt"></div>
    <div class="" :class="{ 'text-center justify-self-center w-3/5': expanded, 'text-left': !expanded }">
      <p
        id="player__podcast-title"
        class="truncate text-primary-dt uppercase"
        :class="{
          'font-extrabold text-xl': expanded,
          'font-bold text-base': !expanded,
        }"
      >{{ podcastTitle }}</p>
      <p
        id="player__episode-title"
        class="truncate"
        :class="{
          'text-sm': !expanded
        }"
      >{{ episodeTitle }}</p>
    </div>
    <div v-show="expanded" v-html="moreVerticalIcon" class="self-center text-primary-dt"></div>
  </div>

  <audio ref="audioElement" :src="audioSrc" preload="auto"></audio>

  <div
    :class="{
      'flex flex-col gap-2 justify-items-start mx-6 my-6': expanded,
      'absolute w-full top-0': !expanded
    }"
  >
    <div class="flex-grow bg-gradient-to-r from-gray-500 to-gray-700  h-1 shadow-inner flex" :class="{ 'rounded-md': expanded }">
      <div class="h-full w-0 shadow-inner" :class="{ 'rounded-md': expanded }" :style="'background-color: #14B8A6; width: ' + calculatedProgress  + '%;'"></div>
      <div v-show="expanded" class="h-4 w-4 rounded-full border border-black relative -left-2" style="background-color: #14B8A6; top: -6px;"></div>
    </div>
    <div v-show="expanded" class="flex flex-row font-bold text-sm bg-transparent mx-2 justify-between">
      <p>{{ currentTimeStr }}</p>
      <p>-{{ remainingTimeStr }}</p>
    </div>
  </div>

  <div
    id="player__buttons"
    class="bg-transparent transition-colors duration-500 text-primary-dt"
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
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="0.8" stroke-linecap="round" stroke-linejoin="round" class="w-20 h-20 md:w-28 md:h-28"><circle cx="12" cy="12" r="10"></circle><polygon points="10 8 16 12 10 16 10 8" fill="black" stroke="black"></polygon></svg>
        </div>
        <div v-else>
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="0.8" stroke-linecap="round" stroke-linejoin="round" class="w-20 h-20 md:w-28 md:h-28"><circle cx="12" cy="12" r="10"></circle><line x1="10" y1="15" x2="10" y2="9" stroke="black" stroke-width="2"></line><line x1="14" y1="15" x2="14" y2="9" stroke="black" stroke-width="2"></line></svg>
        </div>
      </div>
      <div v-else>
        <div v-if="!playing">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-9 h-9"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>
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
  <div id="player__episode-description" v-if="episodeDescription !== ''" v-show="expanded" class="m-4 p-2 rounded-xl text-primary-dt bg-secondary-dt">
    {{ episodeDescription }}
  </div>
</div>
</template>

<script lang='ts'>
import {
  defineComponent,
  ref,
  computed,
  onMounted,
  onBeforeUnmount,
} from 'vue';
import feather from 'feather-icons';
import { PlayerAPI, SubscriptionsAPI } from '@/api';
import { playerEventBus, PlayerEvents } from '@/events/player';
import { Episode } from '@/api/types';
import defaultArtwork from '@/assets/resources/default_artwork.svg';

export default defineComponent({
  props: {
    expanded: {
      type: Boolean,
      required: false,
      default: false,
    },
  },
  emits: ['open-request', 'close-request'],
  setup(_, context) {
    const playerAPI = new PlayerAPI();
    const subsAPI = new SubscriptionsAPI();

    const currentEpisode = ref<Episode | null>(null);

    const playing = ref(false);
    const audioElement = ref<HTMLAudioElement | null>(null);
    const currentTime = ref(0);
    const remainingTime = ref(0);
    const currentTimeStr = ref('00:00');
    const remainingTimeStr = ref('00:00');
    const duration = ref(0);

    const audioSrc = ref('');
    const artworkSrc = ref('');
    const podcastTitle = ref('');
    const episodeTitle = ref('');
    const episodeDescription = ref('');

    const playEpisode = async (episode: Episode) => {
      currentEpisode.value = episode;

      audioSrc.value = episode.enclosureURL;
      episodeTitle.value = episode.title;
      episodeDescription.value = episode.description;

      const podcast = await subsAPI.getPodcastDetails(episode.parentPodcastID);

      if (podcast === undefined) {
        throw new Error(`Unable to find the podcast with ID '${episode.ID}'`);
      }

      artworkSrc.value = podcast.imageURL;
      podcastTitle.value = podcast.title;

      if (audioElement.value == null) {
        throw new Error('Unable to load the audio because the referenced audio element is null');
      }

      audioElement.value.load();
      audioElement.value.currentTime = episode.currentProgress;
    };

    // Here we handle the request to play episodes
    playerEventBus.on(PlayerEvents.PLAY_REQUEST, async (e) => {
      const episode = e as Episode;

      await playEpisode(episode);
    });

    const rotateCwIcon = computed(() => feather.icons['rotate-cw'].toSvg({ 'stroke-width': 1.5, class: 'w-8 h-8 md:w-12 md:h-12' }));
    const rotateCcwIcon = computed(() => feather.icons['rotate-ccw'].toSvg({ 'stroke-width': 1.5, class: 'w-8 h-8 md:w-12 md:h-12' }));
    const skipBackIcon = computed(() => feather.icons['skip-back'].toSvg({ 'stroke-width': 1.5, class: 'w-6 h-6 md:w-10 md:h-10' }));
    const skipForwardIcon = computed(() => feather.icons['skip-forward'].toSvg({ 'stroke-width': 1.5, class: 'w-6 h-6 md:w-10 md:h-10' }));
    const share2Icon = computed(() => feather.icons['share-2'].toSvg({ class: 'mx-6 md:mx-14' }));
    const moreVerticalIcon = computed(() => feather.icons['more-vertical'].toSvg({ class: 'mx-6 md:mx-12' }));

    const playPause = () => {
      if (audioElement.value == null) {
        throw new Error('audioElement null, unable to resume/pause the reproduction');
      }

      if (!playing.value) {
        audioElement.value.play();
        playerEventBus.emit(PlayerEvents.PLAYBACK_STATUS_CHANGE, 'play');
      } else {
        audioElement.value.pause();
        playerEventBus.emit(PlayerEvents.PLAYBACK_STATUS_CHANGE, 'pause');
      }
    };

    const skipBackward = (secs: number) => {
      if (audioElement.value == null) {
        throw new Error('audioElement null, unable to skip backward');
      }

      if (currentTime.value <= secs) {
        audioElement.value.currentTime = 0;
        currentTime.value = 0;
      } else {
        audioElement.value.currentTime -= secs;
        currentTime.value -= secs;
      }
    };

    const skipForward = (secs: number) => {
      if (audioElement.value == null) {
        throw new Error('audioElement null, unable to skip forward');
      }

      if ((currentTime.value + secs) >= duration.value) {
        audioElement.value.currentTime = duration.value;
        currentTime.value = duration.value;
      } else {
        audioElement.value.currentTime += secs;
        currentTime.value += secs;
      }
    };

    const secsToMMSS = (secs: number) => {
      const minutes = Math.floor(secs / 60);
      const seconds = Math.floor(secs - (minutes * 60));

      function c(n: number): string {
        return n < 10 ? `0${n}` : `${n}`;
      }

      return `${c(minutes)}:${c(seconds)}`;
    };

    const calculatedProgress = computed((): number => (currentTime.value * 100) / duration.value);

    const updateRemaining = () => {
      if (audioElement.value == null) {
        throw new Error('audioElement null, unable to update the remaining time of the playback');
      }

      remainingTime.value = Math.floor(duration.value) - Math.floor(currentTime.value);
      remainingTimeStr.value = secsToMMSS(remainingTime.value);
    };

    const updateDuration = () => {
      if (audioElement.value == null) {
        throw new Error('audioElement null, unable to update the duration of the playback');
      }

      duration.value = audioElement.value.duration;
      updateRemaining();
    };

    const updateCurrentAndRemaining = () => {
      if (audioElement.value == null) {
        throw new Error('audioElement null, unable to update the duration and remaining time of the playback');
      }

      if (currentEpisode.value == null) {
        throw new Error('The variable that contains the episode that is being played shouldn\'t be null');
      }

      playerAPI.updateEpisodeProgress(currentEpisode.value.parentPodcastID, currentEpisode.value.ID, currentTime.value)
        .catch((err) => {
          throw err;
        });

      currentTime.value = audioElement.value.currentTime;
      currentTimeStr.value = secsToMMSS(currentTime.value);
      updateRemaining();

      playerEventBus.emit(PlayerEvents.PROGRESS_CHANGE, currentTime.value);
    };

    const setPlaying = () => { playing.value = true; };

    const setPaused = () => { playing.value = false; };

    const emitOpenEvent = () => {
      context.emit('open-request');
    };

    const emitCloseEvent = () => {
      context.emit('close-request');
    };

    const onEnded = (_event: Event) => {
      artworkSrc.value = '';
      podcastTitle.value = '';
      episodeTitle.value = '';
      audioSrc.value = '';
      episodeDescription.value = '';
      currentEpisode.value = null;

      // TODO send a request to set the episode as played.

      if (audioElement.value == null) {
        throw new Error('audioElement null, unable to reset the audio src');
      }

      audioElement.value.load();
      playerEventBus.emit(PlayerEvents.PLAYBACK_END);
    };

    const onError = (err: ErrorEvent) => {
      playerEventBus.emit(PlayerEvents.ERROR, err.message);
      throw err;
    };

    onMounted(async () => {
      if (!audioElement.value) {
        return;
      }

      audioElement.value.addEventListener('durationchange', updateDuration);
      audioElement.value.addEventListener('timeupdate', updateCurrentAndRemaining);
      audioElement.value.addEventListener('play', setPlaying);
      audioElement.value.addEventListener('pause', setPaused);
      audioElement.value.addEventListener('ended', onEnded);
      audioElement.value.addEventListener('error', onError);

      // Restore the status of the player
      const playbackInfo = await playerAPI.getPlayerPlaybackInfo();

      // TODO use the endpoint that returns details about a specific episode (depends of https://github.com/LinCast-Team/LinCast/issues/182)
      const podcastEpisodes = await subsAPI.getEpisodes(playbackInfo.podcastID);
      const episode = podcastEpisodes.find((e) => e.ID === playbackInfo.episodeID);

      if (episode == null) {
        throw new Error('Episode to be played not found');
      }

      await playEpisode(episode);
    });

    onBeforeUnmount(() => {
      if (!audioElement.value) {
        return;
      }

      audioElement.value.removeEventListener('durationchange', updateDuration);
      audioElement.value.removeEventListener('timeupdate', updateCurrentAndRemaining);
      audioElement.value.removeEventListener('play', setPlaying);
      audioElement.value.removeEventListener('pause', setPaused);
      audioElement.value.removeEventListener('ended', onEnded);
      audioElement.value.removeEventListener('error', onError);
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
      emitCloseEvent,

      // Player elements
      audioSrc,
      artworkSrc,
      podcastTitle,
      episodeTitle,
      episodeDescription,
      defaultArtwork,
    };
  },
});
</script>

<style lang="scss">
</style>
