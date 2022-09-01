import {
  Podcast,
  Episode,
} from './types';
import APIBase from './api-base';

class SubscriptionsAPI extends APIBase {
  static async getSubscriptions(): Promise<Array<Podcast>> {
    const response = await fetch(`${SubscriptionsAPI.BASE_PATH}/user/subscriptions`, {
      method: 'GET',
    });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Request failed with status code ${response.status}: ${body}`);
    }

    const data: Array<Podcast> = await response.json();

    return data;
  }

  static async getPodcastDetails(podcastID: number): Promise<Podcast> {
    const response = await fetch(`${this.BASE_PATH}/podcasts/${podcastID}`, {
      method: 'GET',
    });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Request failed with status code ${response.status}: ${body}`);
    }

    const data: Podcast = await response.json();

    return data;
  }

  static async getEpisodes(podcastID: number): Promise<Array<Episode>> {
    const response = await fetch(`${this.BASE_PATH}/podcasts/${podcastID}/episodes`, {
      method: 'GET',
    });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Request failed with status code ${response.status}: ${body}`);
    }

    const data: Array<Episode> = await response.json();

    return data;
  }

  static async subscribe(feedURL: string) {
    const response = await fetch(`${this.BASE_PATH}/podcasts/subscribe`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ url: feedURL }),
    });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Request failed with status code ${response.status}: ${body}`);
    }
  }

  static async unsubscribe(podcastID: number) {
    const response = await fetch(`${this.BASE_PATH}/podcasts/unsubscribe?id=${podcastID}`, {
      method: 'PUT',
    });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Request failed with status code ${response.status}: ${body}`);
    }
  }

  static async getLatestSubscriptionsEpisodes(from: string, to: string): Promise<Array<Episode>> {
    const response = await fetch(`${this.BASE_PATH}/podcasts/latest_eps?from=${from}&to=${to}`, {
      method: 'GET',
    });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Request failed with status code ${response.status}: ${body}`);
    }

    const data: Array<Episode> = await response.json();

    return data;
  }
}

export default SubscriptionsAPI;
