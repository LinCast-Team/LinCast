export interface PlaybackInfo {
  episodeID: string;
  podcastID: number;
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
