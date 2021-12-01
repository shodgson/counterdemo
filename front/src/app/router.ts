import { createRouter, createWebHistory } from "vue-router";
//import SignIn from "@/views/SignIn.vue";
import SignUp from "@/views/SignUp.vue";
//import AccountSetup from "@/views/AccountSetup.vue";
//import Dashboard from "@/views/Dashboard.vue";
import store from "@/store";

const routes = [
  { path: "/", name: "default", component: SignUp, meta: { requiresAuth: true } },
  //{ path: "/signin", name: "signIn", component: SignIn },
  { path: "/signup", name: "signUp", component: SignUp },
  { path: "/signup2", name: "profile", component: SignUp },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.VITE_APP_PATH),
  routes: routes,
  //base: import.meta.env.VITE_APP_PATH,
});

router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth && !store.getters["cognito/isAuthenticated"])
    next({ name: "signUp" });
  else next();
});

export default router;
