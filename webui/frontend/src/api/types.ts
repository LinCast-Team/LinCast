export interface Podcast {
  ID: number;
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
  ID: number;
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
  currentProgress: number;
}

export interface PlaybackInfo {
  podcastID: number;
  episodeID: number;
}

export type Queue = Array<QueueEpisode>;

export interface QueueEpisode {
  id: number;
  podcastID: number;
  episodeID: string;
  positon: number;
}
