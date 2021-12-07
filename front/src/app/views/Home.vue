<template>
  <div>
    <div class="text-xl mt-6">Counter</div>
    <div v-if="loading">
      <div class="mt-36 flex justify-center items-center">
        <div
          class="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-neutral"
        ></div>
      </div>
    </div>
    <div v-if="!loading" class="">
      <div class="mx-auto p-8 font-bold text-4xl rounded-full shadow bg-accent text-accent-content w-32 mt-8">
        {{ count }}
      </div>
      <div>
        <button
          @click="increment"
          class="btn btn-sm mt-8"
          :class="{ loading: btnIncLoading }"
        >
          Increment
        </button>
        </div><div>
        <button
          @click="reset"
          class="btn btn-sm mt-4"
          :class="{ loading: btnResetLoading }"
        >
          Reset
        </button>
      </div>
      <div v-if="errorText" class="w-1/2 mt-8 mx-auto alert alert-error">
  <div class="flex-1">
    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="w-6 h-6 mx-2 stroke-current">    
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"></path>                      
    </svg> 
    <label>{{ errorText }}</label>
  </div>
</div>

    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import store from "@/store";
import { mapGetters } from "vuex";
import countApi from "@/api";
export default defineComponent({
  components: {},
  data() {
    return {
      loading: true,
      count: -1,
      btnResetLoading: false,
      btnIncLoading: false,
      errorText: "",
    };
  },
  computed: {
    ...mapGetters("cognito", ["username"]),
    ...mapGetters("account", ["accountActive"]),
  },
  created() {
    countApi
      .getUser(this.username)
      .then((r) => this.count = r.count)
      .catch((e) => this.errorText = e)
      .finally(() => this.loading = false)
    /*loadProfile().then((r) => {
      this.series = r.series;
      store.dispatch("account/setActiveStatus", r.active);
      this.loading = false;
    });
    const paymentResult = this.$route.query.payment;
    if (paymentResult) {
      if (paymentResult == "failure") {
        this.$toast.open({
          type: "failure",
          message: "Unable to upgrade account",
        });
      } else if (paymentResult == "success") {
        this.$toast.open({
          type: "success",
          message: "Successfully upgraded account!",
        });
      }
    }
    */
  },
  methods: {
    clearError() {
      this.errorText = ""
    },
    increment() {
      this.clearError()
      this.btnIncLoading = true
      this.count++
      countApi
        .add(this.username, 1)
        .then(r => this.count = r.count)
      .catch((e) => {
        this.errorText = e
        this.count--;
      })
        .finally(() => this.btnIncLoading = false);
    },
    reset() {
      this.clearError()
      this.btnResetLoading = true;
      const tmpCount = this.count
      this.count = 0;
      countApi
        .reset(this.username)
        .then(r => this.count = r.count)
      .catch((e) => {
        this.errorText = e
        this.count = tmpCount
      })
        .finally(() => this.btnResetLoading = false);
    },
  },
});
</script>
<style></style>
