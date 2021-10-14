import { ActionTree, ActionContext } from 'vuex';
import { State } from './state';
import { Mutations } from './mutations';
import MutationTypes from './mutation-types';
import ActionTypes from './action-types';
import {
  Podcast,
  Episode,
  getUserSubscriptions,
  getLatestEpisodes,
} from './api/subscriptions';
import {
  PlaybackInfo,
  getPlaybackInfo,
  updatePlaybackInfo,
  getEpisodeProgress,
  updateEpisodeProgress,
} from './api/player';

type AugmentedActionContext = {
  commit<K extends keyof Mutations>(
    key: K,
    payload: Parameters<Mutations[K]>[1]
  ): ReturnType<Mutations[K]>;
} & Omit<ActionContext<State, State>, 'commit'>

export interface Actions {
  [ActionTypes.GET_SUBSCRIPTIONS]({ commit }: AugmentedActionContext): Promise<Array<Podcast>>;
  [ActionTypes.GET_LATEST_EPISODES]({ commit }: AugmentedActionContext, payload: { from: string; to: string }): Promise<Array<Episode>>;
  [ActionTypes.GET_PLAYBACK_INFO]({ commit }: AugmentedActionContext): Promise<PlaybackInfo>;
  [ActionTypes.SET_PLAYBACK_INFO]({ commit }: AugmentedActionContext, payload: PlaybackInfo): Promise<void>;
  [ActionTypes.GET_PROGRESS]({ commit, state }: AugmentedActionContext): Promise<number>;
  [ActionTypes.SET_PROGRESS]({ commit, state }: AugmentedActionContext, payload: number): Promise<void>;
}

export const actions: ActionTree<State, State> & Actions = {
  async [ActionTypes.GET_SUBSCRIPTIONS]({ commit }) {
    const data = await getUserSubscriptions();
    commit(MutationTypes.SET_SUBSCRIPTIONS, data);

    return data;
  },
  async [ActionTypes.GET_LATEST_EPISODES]({ commit }, payload) {
    const data = await getLatestEpisodes(payload.from, payload.to);
    commit(MutationTypes.SET_LATEST_EPISODES, data);

    return data;
  },
  async [ActionTypes.GET_PLAYBACK_INFO]({ commit }) {
    const data = await getPlaybackInfo();
    commit(MutationTypes.SET_PLAYBACK_INFO, data);

    return data;
  },
  async [ActionTypes.SET_PLAYBACK_INFO]({ commit }, payload) {
    await updatePlaybackInfo(payload);
    commit(MutationTypes.SET_PLAYBACK_INFO, payload);
  },
  async [ActionTypes.GET_PROGRESS]({ commit, state }) {
    const progress = await getEpisodeProgress(state.playbackInfo.podcastID, state.playbackInfo.episodeID);
    commit(MutationTypes.SET_PROGRESS, progress);

    return progress;
  },
  async [ActionTypes.SET_PROGRESS]({ commit, state }, payload) {
    commit(MutationTypes.SET_PROGRESS, payload);
    await updateEpisodeProgress(state.playbackInfo.podcastID, state.playbackInfo.episodeID, payload);
  },
};
