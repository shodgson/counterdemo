<template>
  <div class="max-w-xl mx-auto">
    <div class="text-xl mt-6">Profile</div>

    <div class="text-left">
      <div class="flex justify-between">
        <div class="font-bold">Email</div>
        <div v-if="ready">{{ email }}</div>
      </div>

      <div class="flex justify-between mt-4">
        <div class="font-bold">User ID</div>
        <div v-if="ready">{{ username }}</div>
      </div>

      <div class="flex justify-between mt-4">
        <div class="font-bold">Current plan</div>
        <div v-if="accountActive && ready" class="flex">
          <div class="text-right">Premium (US$20/month)</div>
          <button
            class="btn btn-outline btn-xs ml-2"
            :class="{ loading: loadActivation }"
            @click="cancel"
          >
            Cancel
          </button>
        </div>
        <div v-if="!accountActive && ready" class="flex">
          <div class="text-right">Free</div>
          <button
            class="btn btn-outline btn-xs ml-2"
            :class="{ loading: loadActivation }"
            @click="upgrade"
          >
            Upgrade
          </button>
        </div>
      </div>
    </div>
    <div>
      <button
        @click="signOut"
        class="btn btn-sm mt-8"
        :class="{ loading: signingOut }"
      >
        Sign out
      </button>
    </div>
    <div class="mt-5 text-gray-500 text-sm" v-if="!accountActive && ready">
      Demo: you can upgrade with a <a class="link" href="https://stripe.com/docs/testing#cards">Stripe dummy credit card</a>

      </div>

  </div>
</template>
<script lang="ts">
import store from "@/store";
import { mapGetters } from "vuex";
import { defineComponent } from "vue";
import countApi from "@/api";
export default defineComponent({
  components: {},
  data() {
    return {
      signingOut: false,
      loadActivation: false,
      email: "",
      accountActive: false,
      ready: false,
    };
  },
  computed: {
    ...mapGetters("cognito", ["username"]),
  },
  mounted() {
    countApi
      .getUser(this.username)
      .then((r) => {
        this.email = r.name
        this.accountActive = r.premium || false
        this.ready = true
      })

  },
  methods: {
    signOut: function () {
      this.signingOut = true;
      store.dispatch("cognito/signOut");
      store.dispatch("account/setActiveStatus", null);
      this.signingOut = false;
      this.$router.push({ name: "signIn" });
    },
    upgrade: function () {
      this.loadActivation = true;
      countApi.paymentUrl().then((r) => {
        this.loadActivation = false;
        window.location.href = r.URL;
      });
    },
    cancel: function () {
      this.loadActivation = true;
      countApi
        .cancelSubscription()
        .then(() => this.accountActive = false)
        .catch((e) => console.error(e))
        .finally(() => this.loadActivation = false)

      /*
      cancelAccount()
        .then((r) => {
          loadProfile().then((p) => {
            store.dispatch("account/setActiveStatus", p.active);
          });
          this.$toast.open({
            message: "Subscription cancelled",
            type: "success",
          });
          this.loadActivation = false;
        })
        .catch((e) => {
          this.loadActivation = false;
          e.then((p) => {
            console.error(p.message)
          });
        });
        */
    },
  },
});
</script>
<style scoped></style>
