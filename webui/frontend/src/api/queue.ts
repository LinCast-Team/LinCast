export type Queue = Array<QueueEpisode>;

export interface QueueEpisode {
  ID: number;
  podcastID: number;
  episodeID: string;
  positon: number;
}

export const getQueue = async (): Promise<Queue> => {
  const response = await fetch('/api/v0/player/queue', {
    method: 'GET'
  });

  if (!response.ok) {
    throw new Error(`Request failed with status code ${response.status}`);
  }

  const data: Queue = await response.json();

  return data;
}

export const clearQueue = async () => {
  const response = await fetch('/api/v0/player/queue', {
    method: 'DELETE'
  });

  if (!response.ok) {
    throw new Error(`Request failed with status code ${response.status}`);
  }
}

export const replaceQueue = async (newQueue: Queue) => {
  const response = await fetch('/api/v0/player/queue', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(newQueue),
  });

  if (!response.ok) {
    throw new Error(`Request failed with status code ${response.status}`);
  }
}

export const addToQueue = async (episode: QueueEpisode, append: boolean) => {
  const response = await fetch(`/api/v0/player/queue/add?append=${append}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(episode),
  });

  if (!response.ok) {
    throw new Error(`Request failed with status code ${response.status}`);
  }
}

export const removeFromQueue = async (id: number) => {
  const response = await fetch(`/api/v0/player/queue/remove?id=${id}`, {
    method: 'DELETE',
  });

  if (!response.ok) {
    throw new Error(`Request failed with status code ${response.status}`);
  }
}
