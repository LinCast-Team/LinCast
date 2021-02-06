<template>
    <div class="p-3 w-64 bg-gray-800 flex rounded-xl border border-green-700">
      <div v-html="searchIcon"></div>
      <input
        class="w-auto bg-transparent border-0 placeholder-gray-50 text-gray-50 focus:placeholder-gray-500 focus:outline-none"
        type="text"
        role="search"
        name="Search a podcast"
        placeholder="Search"
        aria-placeholder="Search"
        v-model.lazy="data"
      >
    </div>
</template>

<script>
import { computed, ref, watch } from 'vue';
import feather from 'feather-icons';

export default {
  emits: ['search-input'],
  setup(_, context) {
    const searchIcon = computed(() => feather.icons.search.toSvg({ 'stroke-width': 1.5, class: 'w-5 h-5 md:w-12 md:h-12 text-gray-50 mr-2' }));

    const data = ref('');

    watch(data, (_oldVal, newVal) => {
      context.emit('search-input', newVal);
    });

    return {
      data,
      searchIcon,
    };
  },
};
</script>
