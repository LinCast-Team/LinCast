/* eslint-disable class-methods-use-this */
import {
  Module, VuexModule, MutationAction, getModule,
} from 'vuex-module-decorators';
import store from '@/store';
import { CurrentProgress, getCurrentProgress, sendCurrentProgress } from '../api/player';
import {
  Queue, QueueEpisode, getQueue, addToQueue,
} from '../api/queue';

@Module({
  namespaced: true,
  name: 'player',
  dynamic: true,
  store,
})

class Player extends VuexModule {
  currentProgress: CurrentProgress = {} as CurrentProgress;

  queueContent: Queue = {} as Queue;

  get progress() {
    return this.currentProgress;
  }

  get queue() {
    return this.queueContent;
  }

  @MutationAction
  async fetchProgress() {
    const currentProgress = await getCurrentProgress();

    return {
      currentProgress,
    };
  }

  @MutationAction
  async updateProgress(newProgress: CurrentProgress) {
    await sendCurrentProgress(newProgress);

    return {
      currentProgress: newProgress,
    };
  }

  @MutationAction
  async fetchQueue() {
    const queueContent = await getQueue();

    return {
      queueContent,
    };
  }

  @MutationAction
  async addToQueue(ep: QueueEpisode, append: boolean) {
    await addToQueue(ep, append);

    if (append) {
      this.queue.push(ep);
    } else {
      this.queue.unshift(ep);
    }

    return {};
  }
}

export default getModule(Player);
