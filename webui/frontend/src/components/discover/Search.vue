<template>
    <div class="p-3 bg-gray-800 flex " :class="{ 'w-11/12 rounded-xl': !searchMode, 'w-full rounded-none': searchMode }">
      <div v-if="!searchMode" v-html="searchIcon"></div>
      <div v-else @click="onFocusChange(false)" v-html="backIcon"></div>

      <input
        class="w-auto bg-transparent border-0 placeholder-gray-100 text-gray-50 focus:placeholder-gray-500 focus:outline-none"
        type="text"
        role="search"
        name="Search a Podcast"
        placeholder="Podcast..."
        aria-placeholder="Search a Podcast"
        v-model.lazy="data"
        @focus="onFocusChange(true)"
      >
    </div>
</template>

<script lang='ts'>
import {
  defineComponent,
  computed,
  ref,
  watch,
} from 'vue';
import feather from 'feather-icons';

export default defineComponent({
  emits: ['search-input', 'search-focus'],
  props: {
    searchMode: {
      required: true,
      type: Boolean,
    },
  },
  setup(_, context) {
    const searchIcon = computed(() => feather.icons.search.toSvg({ 'stroke-width': 1.5, class: 'w-5 h-5 md:w-12 md:h-12 text-gray-100 mr-2' }));
    const backIcon = computed(() => feather.icons['chevron-left'].toSvg({ 'stroke-width': 1.5, class: 'w-5 h-5 md:w-12 md:h-12 text-gray-100 mr-2' }));

    const data = ref('');

    watch(data, (_oldVal, newVal) => {
      context.emit('search-input', newVal);
    });

    const onFocusChange = (val: boolean) => {
      context.emit('search-focus', val);
    };

    return {
      data,
      searchIcon,
      backIcon,
      onFocusChange,
    };
  },
});
</script>
