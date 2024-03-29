import {
  PlaybackInfo,
} from './types';
import APIBase from './api-base';

class PlayerAPI extends APIBase {
  async getPlayerPlaybackInfo(): Promise<PlaybackInfo> {
    const response = await fetch(`${this.BASE_PATH}/player/playback_info`, { method: 'GET' });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Request failed with status code ${response.status}: ${body}`);
    }

    const data: PlaybackInfo = await response.json();

    return data;
  }

  async postPlayerPlaybackInfo(p: PlaybackInfo) {
    const response = await fetch(`${this.BASE_PATH}/player/playback_info`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(p),
    });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Request failed with status code ${response.status}: ${body}`);
    }
  }

  async getEpisodeProgress(podcastID: number, episodeID: number) {
    const response = await fetch(`${this.BASE_PATH}/podcasts/${podcastID}/episodes/${episodeID}/progress`, { method: 'GET' });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Request failed with status code ${response.status}: ${body}`);
    }

    const { progress } = await response.json() as { progress: number };

    return progress;
  }

  async updateEpisodeProgress(podcastID: number, episodeID: number, newProgress: number) {
    const response = await fetch(`${this.BASE_PATH}/podcasts/${podcastID}/episodes/${episodeID}/progress`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ progress: Math.round(newProgress) }),
    });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Request failed with status code ${response.status}: ${body}`);
    }
  }
}

export default PlayerAPI;
