import {
  Podcast,
  Episode,
} from './types';

class SubscriptionsAPI extends APIBase {
  async getSubscriptions(): Promise<Array<Podcast>> {
    const response = await fetch(`${this.BASE_PATH}/user/subscriptions`, {
      method: 'GET',
    });

    if (!response.ok) {
      throw new Error(`Request failed with status code ${response.status}`);
    }

    const data: Array<Podcast> = await response.json();

    return data;
  }

  async getPodcastDetails(podcastID: number): Promise<Podcast> {
    const response = await fetch(`${this.BASE_PATH}/podcasts/${podcastID}/details`, {
      method: 'GET',
    });

    if (!response.ok) {
      throw new Error(`Request failed with status code ${response.status}`);
    }

    const data: Podcast = await response.json();

    return data;
  }

  async getEpisodes(podcastID: number): Promise<Array<Episode>> {
    const response = await fetch(`${this.BASE_PATH}/podcasts/${podcastID}/episodes`, {
      method: 'GET',
    });

    if (!response.ok) {
      throw new Error(`Request failed with status code ${response.status}`);
    }

    const data: Array<Episode> = await response.json();

    return data;
  }

  async subscribe(feedURL: string) {
    const response = await fetch(`${this.BASE_PATH}/podcasts/subscribe`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ url: feedURL }),
    });

    if (!response.ok) {
      throw new Error(`Request failed with status code ${response.status}`);
    }
  }

  async unsubscribe(podcastID: number) {
    const response = await fetch(`${this.BASE_PATH}/podcasts/unsubscribe?id=${podcastID}`, {
      method: 'PUT',
    });

    if (!response.ok) {
      throw new Error(`Request failed with status code ${response.status}`);
    }
  }

  async getLatestSubscriptionsEpisodes(from: string, to: string): Promise<Array<Episode>> {
    const response = await fetch(`${this.BASE_PATH}/podcasts/latest_eps?from=${from}&to=${to}`, {
      method: 'GET',
    });

    if (!response.ok) {
      throw new Error(`Request failed with status code ${response.status}`);
    }

    const data: Array<Episode> = await response.json();

    return data;
  }
}

export default SubscriptionsAPI;
