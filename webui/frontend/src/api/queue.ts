import {
  Queue,
  QueueEpisode,
} from './types';
import APIBase from './api-base';

class PlayerQueueAPI extends APIBase {
  async getQueue(): Promise<Queue> {
    const response = await fetch(`${this.BASE_PATH}/player/queue`, {
      method: 'GET',
    });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Request failed with status code ${response.status}: ${body}`);
    }

    const data: Queue = await response.json();

    return data;
  }

  async clearQueue() {
    const response = await fetch(`${this.BASE_PATH}/player/queue`, {
      method: 'DELETE',
    });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Request failed with status code ${response.status}: ${body}`);
    }
  }

  async overwriteQueue(newQueue: Queue) {
    const response = await fetch(`${this.BASE_PATH}/player/queue`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(newQueue),
    });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Request failed with status code ${response.status}: ${body}`);
    }
  }

  async addToQueue(episode: QueueEpisode, append: boolean) {
    const response = await fetch(`${this.BASE_PATH}/player/queue/add?append=${append}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(episode),
    });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Request failed with status code ${response.status}: ${body}`);
    }
  }

  async removeFromQueue(id: number) {
    const response = await fetch(`${this.BASE_PATH}/player/queue/remove?id=${id}`, {
      method: 'DELETE',
    });

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`Request failed with status code ${response.status}: ${body}`);
    }
  }
}

export default PlayerQueueAPI;
