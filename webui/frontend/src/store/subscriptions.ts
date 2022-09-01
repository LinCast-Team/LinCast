import { defineStore } from 'pinia';
import dayjs from 'dayjs';
import { Podcast, Episode } from '@/api/types';
import { SubscriptionsAPI } from '@/api';

/* eslint-disable import/prefer-default-export */
export const useSubscriptionsStore = defineStore('subscriptions', {
  state: () => ({
    podcasts: new Array<Podcast>(),
    episodes: new Array<Episode>(),
    subscriptionsFrom: dayjs().subtract(30, 'day').format('YYYY-MM-DD'),
    subscriptionsTo: dayjs().format('YYYY-MM-DD'),
  }),
  getters: {},
  actions: {
    async fetchPodcasts(): Promise<void> {
      const subs = await SubscriptionsAPI.getSubscriptions();
      this.podcasts = subs;
    },
    async fetchEpisodes(): Promise<void> {
      const eps = await SubscriptionsAPI.getLatestSubscriptionsEpisodes(this.subscriptionsFrom, this.subscriptionsTo);
      this.episodes = eps;
    },
  },
});
