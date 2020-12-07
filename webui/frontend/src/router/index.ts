import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
import Home from '../views/Home.vue';

const DiscoverView = () => import('../views/Discover.vue');
const LibraryView = () => import('../views/Library.vue');

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: Home,
  },
  {
    path: '/discover',
    name: 'Discover',
    component: DiscoverView,
  },
  {
    path: '/library',
    name: 'Library',
    component: LibraryView,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
