import EpisodeItem from '../components/library/EpisodeItem.vue';

export default {
  title: 'EpisodeItem',
  component: EpisodeItem,
};

const Template = (args, { argTypes }) => ({
  components: { EpisodeItem },
  props: Object.keys(argTypes),
  template: '<EpisodeItem v-bind="$props"/>',
});

export const FirstStory = Template.bind({});
FirstStory.args = {
  title: 'Hello world!',
  author: 'Martin',
  resume: 'Something really, really boring.',
  duration: '01:12:24',
};
