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
</template>

<script lang="ts">
import { defineComponent } from "vue";
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
