import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';
import '@/assets/css/tailwind.css';

const app = createApp(App);

app.use(router);
app.use(createPinia());

router.isReady().then(() => {
  app.mount('#app');
});
