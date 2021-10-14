import { MutationTree } from 'vuex';
import MutationTypes from './mutation-types';
import { State } from './state';
import { Podcast, Episode } from './api/subscriptions';
import { PlaybackInfo } from './api/player';

export type Mutations<S = State> = {
  [MutationTypes.SET_SUBSCRIPTIONS](state: S, payload: Array<Podcast>): void;
  [MutationTypes.SET_LATEST_EPISODES](state: S, payload: Array<Episode>): void;
  [MutationTypes.SET_PLAYBACK_INFO](state: S, payload: PlaybackInfo): void;
}

export const mutations: MutationTree<State> & Mutations = {
  [MutationTypes.SET_SUBSCRIPTIONS](state, payload: Array<Podcast>) {
    state.userPodcasts = payload;
  },
  [MutationTypes.SET_LATEST_EPISODES](state, payload: Array<Episode>) {
    state.latestEpisodes = payload;
  },
  [MutationTypes.SET_PLAYBACK_INFO](state, payload: PlaybackInfo) {
    state.playbackInfo = payload;
  },
};
