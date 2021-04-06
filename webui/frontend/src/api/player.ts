export interface CurrentProgress {
  progress: number;
  episodeID: string;
  podcastID: number;
}

export const getCurrentProgress = async (): Promise<CurrentProgress> => {
  const response = await fetch('/api/v0/player/progress', { method: 'GET' });

  if (!response.ok) {
    throw new Error(`Request failed with status code ${response.status}`);
  }

  const data = await response.json();

  const progress: CurrentProgress = {
    progress: data.progress,
    episodeID: data.episodeID,
    podcastID: data.podcastID,
  };

  return progress;
};

export const sendCurrentProgress = async (p: CurrentProgress) => {
  const request = await fetch('/api/v0/player/progress', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(p),
  });

  if (!request.ok) {
    throw new Error(`Request failed with status code ${request.status}`);
  }
};
