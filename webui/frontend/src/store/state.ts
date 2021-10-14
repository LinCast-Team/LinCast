import { Podcast, Episode } from './api/subscriptions';
import { PlaybackInfo } from './api/player';

export const state = {
  userPodcasts: new Array<Podcast>(),
  latestEpisodes: new Array<Episode>(),
  playbackInfo: {} as PlaybackInfo,
  playerProgress: 0,
};

export type State = typeof state;
