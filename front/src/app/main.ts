import { createApp } from 'vue'
import App from './App.vue'
import store from "@/store"
import router from "@/router";
import "@/index.css"

const app = createApp(App)
app.use(router)
app.use(store)
store
  .dispatch("cognito/getCurrentUser")
  .catch((e) => {
    console.info("No user:", e.message);
  })
  .finally(() => {
    app.mount("#app");
  });
