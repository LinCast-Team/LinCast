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
  title: '',
  author: '',
  resume: '',
  duration: '',
};
