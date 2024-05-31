import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import Toast from "vue-toastification";
import { type PluginOptions } from "vue-toastification";
import "vue-toastification/dist/index.css";

const app = createApp(App)

const toastOptions: PluginOptions = {
    // You can set your default options here
};

app.use(Toast, toastOptions);
app.use(createPinia())
app.use(router)

app.mount('#app')
