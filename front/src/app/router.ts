import { createRouter, createWebHistory } from "vue-router";
import SignIn from "@/views/SignIn.vue";
import SignUp from "@/views/SignUp.vue";
//import AccountSetup from "@/views/AccountSetup.vue";
import Home from "@/views/Home.vue";
import Profile from "@/views/Profile.vue";
import store from "@/store";

const routes = [
  { path: "/", name: "home", component: Home, meta: { requiresAuth: true } },
  { path: "/signin", name: "signIn", component: SignIn },
  { path: "/signup", name: "signUp", component: SignUp },
  { path: "/profile", name: "profile", component: Profile },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.VITE_APP_PATH as string),
  routes: routes,
  //base: import.meta.env.VITE_APP_PATH,
});

router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth && !store.getters["cognito/isAuthenticated"])
    next({ name: "signUp" });
  else next();
});

export default router;
