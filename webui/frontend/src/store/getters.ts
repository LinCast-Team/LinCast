import { GetterTree } from 'vuex';
import { State } from './state';
import { Podcast, Episode } from './api/subscriptions';
import { PlaybackInfo } from './api/player';

export type Getters = {
  userPodcasts(state: State): Array<Podcast>;
  latestEpisodes(state: State): Array<Episode>;
  playbackInfo(state: State): PlaybackInfo;
  playerProgress(state: State): number;
}

export const getters: GetterTree<State, State> & Getters = {
  userPodcasts: (state) => state.userPodcasts,
  latestEpisodes: (state) => state.latestEpisodes,
  playbackInfo: (state) => state.playbackInfo,
  playerProgress: (state) => state.playerProgress,
};
