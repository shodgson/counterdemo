<template>
  <div
    class="navbar m-2 p-2 shadow-lg bg-neutral text-neutral-content rounded-box"
  >
    <div class="flex-1 px-2 mx-2">
      <router-link :to="{ name: 'home' }">
        <button class="btn">Counter demo</button>
      </router-link>
    </div>
    <div class="flex-none px-2 mx-2">
      <router-link :to="{ name: 'profile' }">
        <button v-if="isAuthenticated" class="btn">Profile</button>
      </router-link>
    </div>
  </div>
  <div v-if="accountActive == false" class="alert alert-warning m-2">
    <div class="flex-1">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        class="w-6 h-6 mx-2 stroke-current"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
        ></path>
      </svg>
      <label
        >Your account is limited to 20 calendar downloads per series. Unlock unlimited downloads for US$5/month.</label
      >
    </div>
    <!--
    <button
      @click="activate"
      class="btn btn-sm btn-outline"
      :class="{ loading: loadActivation }"
    >
      Upgrade account
    </button>-->
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { mapGetters } from "vuex";
//import { activateAccount } from "@/api";
//import store from "@/store";

export default defineComponent({
  name: "NavBar",
  components: {},
  data() {
    return {
      loadActivation: false,
    };
  },
  computed: {
    ...mapGetters("cognito", ["isAuthenticated", "username"]),
    ...mapGetters("account", ["accountActive"]),
  },
  methods: {
   /* activate: function () {
      this.loadActivation = true;
      activateAccount().then((r) => {
        this.loadActivation = false;
        if (r.URL) {
          window.location.href = r.URL;
        } else {
          console.error("Problem fetching Stripe URL");
        }
      });
    },
    */
  },
});
</script>

<style></style>
