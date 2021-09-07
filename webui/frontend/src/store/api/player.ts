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

  const data: CurrentProgress = await response.json();

  return data;
};

export const sendCurrentProgress = async (p: CurrentProgress) => {
  const response = await fetch('/api/v0/player/progress', {
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
