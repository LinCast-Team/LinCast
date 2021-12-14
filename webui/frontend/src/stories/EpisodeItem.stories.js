import EpisodeItem from '../components/library/EpisodeItem.vue';

export default {
  title: 'EpisodeItem',
  component: EpisodeItem,
};

const Template = (args) => ({
  components: { EpisodeItem },
  setup() {
    return { args };
  },
  template: '<EpisodeItem v-bind="args"/>',
});

export const FirstStory = Template.bind({});
FirstStory.args = {
  title: 'Hello world!',
  author: 'Martin',
  resume: 'Something really, really boring.',
  duration: '01:12:24',
};
