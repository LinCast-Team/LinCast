/**
 * Interfaces and functions related with the functionality of the player.
 */

export interface PlaybackInfo {
  podcastID: number;
  episodeID: number;
}

export const getPlaybackInfo = async (): Promise<PlaybackInfo> => {
  const response = await fetch('/api/v0/player/playback_info', { method: 'GET' });

  if (!response.ok) {
    throw new Error(`Request failed with status code ${response.status}`);
  }

  const data: PlaybackInfo = await response.json();

  return data;
};

export const updatePlaybackInfo = async (p: PlaybackInfo) => {
  const response = await fetch('/api/v0/player/playback_info', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(p),
  });

  if (!response.ok) {
    throw new Error(`Request failed with status code ${response.status}`);
  }
};

export const getEpisodeProgress = async (podcastID: number, episodeID: number): Promise<number> => {
  const response = await fetch(`/api/v0/podcasts/${podcastID}/episodes/${episodeID}/progress`, { method: 'GET' });

  if (!response.ok) {
    throw new Error(`Request failed with status code ${response.status}`);
  }

  const { progress } = await response.json() as { progress: number };

  return progress;
};

export const updateEpisodeProgress = async (podcastID: number, episodeID: number, newProgress: number) => {
  const response = await fetch(`/api/v0/podcasts/${podcastID}/episodes/${episodeID}/progress`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ progress: newProgress }),
  });

  if (!response.ok) {
    throw new Error(`Request failed with status code ${response.status}`);
  }
};
