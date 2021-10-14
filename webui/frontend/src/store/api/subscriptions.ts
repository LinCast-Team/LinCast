export interface Podcast {
  id: number;
  subscribed: boolean;
  authorName: string;
  authorEmail: string;
  title: string;
  description: string;
  categories: Array<string>;
  imageURL: string;
  imageTitle: string;
  link: string;
  feedLink: string;
  feedType: string;
  feedVersion: string;
  language: string;
  updated: Date;
  lastCheck: Date;
  added: Date;
}

export interface Episode {
  id: number;
  parentPodcastID: number;
  title: string;
  description: string;
  link: string;
  authorName: string;
  guid: string;
  imageURL: string;
  imageTitle: string;
  categories: Array<string>;
  enclosureURL: string;
  enclosureLength: string;
  enclosureType: string;
  season: string;
  published: Date;
  updated: Date;
  played: boolean;
  currentProgress: Date;
}

export const getUserSubscriptions = async (): Promise<Array<Podcast>> => {
  const response = await fetch('/api/v0/user/subscriptions', {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error(`Request failed with status code ${response.status}`);
  }

  const data: Array<Podcast> = await response.json();

  return data;
};

export const getUserPodcastDetails = async (id: number): Promise<Podcast> => {
  const response = await fetch(`/api/v0/podcasts/${id}/details`, {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error(`Request failed with status code ${response.status}`);
  }

  const data: Podcast = await response.json();

  return data;
};

export const getUserPodcastEpisodes = async (podcastID: number): Promise<Array<Episode>> => {
  const response = await fetch(`/api/v0/podcasts/${podcastID}/episodes`, {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error(`Request failed with status code ${response.status}`);
  }

  const data: Array<Episode> = await response.json();

  return data;
};

export const subscribe = async (feedURL: string) => {
  const response = await fetch('/api/v0/podcasts/subscribe', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ url: feedURL }),
  });

  if (!response.ok) {
    throw new Error(`Request failed with status code ${response.status}`);
  }
};

export const unsubscribe = async (podcastID: number) => {
  const response = await fetch(`/api/v0/podcasts/unsubscribe?id=${podcastID}`, {
    method: 'PUT',
  });

  if (!response.ok) {
    throw new Error(`Request failed with status code ${response.status}`);
  }
};

export const getLatestEpisodes = async (): Promise<Array<Episode>> => {
  const response = await fetch('/api/v0/podcasts/latest_eps', {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error(`Request failed with status code ${response.status}`);
  }

  const data: Array<Episode> = await response.json();

  return data;
};
