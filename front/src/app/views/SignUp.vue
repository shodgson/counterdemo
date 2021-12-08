<template>
  <div>
    <h1 class="text-xl my-6">Sign up</h1>
    <div class="w-full max-w-xs mx-auto">
      <form
        @submit.prevent="signUp"
        class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4"
      >
        <div class="mb-4 form-control">
          <label
            class="label"
          >
            Email
          </label>
          <input
            class="input input-bordered"
            :class="invalidUsername ? 'border-red-500' : ''"
            @change="invalidUsername = false"
            type="email"
            placeholder="Email"
            v-model="email"
          />
          <p v-if="invalidUsername" class="text-red-500 text-xs italic">
            {{ invalidUsernameMessage }}
          </p>
        </div>
        <div class="mb-6 form-control">
          <label
            class="label"
          >
            Password
          </label>
          <input
            class="input input-bordered"
            :class="invalidPass ? 'border-red-500' : ''"
            @change="invalidPass = false"
            type="password"
            placeholder="**********"
            v-model="pass"
          />
          <p v-if="invalidPass" class="text-red-500 text-xs italic">
            {{ invalidPassMessage }}
          </p>
        </div>
        <div class="flex w-full">
          <button
            :disabled="pass.length == 0 || email.length == 0"
            class="btn w-full mx-auto"
            :class="{ loading: loading }"
            type="submit"
          >
            Create account
          </button>
        </div>
        <div class="mt-6">
          Already have an account?<br />
          <router-link :to="{ name: 'signIn' }"
            ><a class="link">Sign in here</a></router-link
          >
        </div>
      </form>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { mapGetters } from "vuex";
import store from "@/store";
import { CognitoError } from "@/types/cognito";
export default  defineComponent({
  components: {},
  data() {
    return {
      email: "",
      pass: "",
      name: "",
      invalidPass: false,
      invalidPassMessage: "",
      invalidUsername: false,
      invalidUsernameMessage: "",
      loading: false,
    };
  },
  mounted() {
  },
  computed: {
    //...mapState({
    //  user: state => state.auth.user,
    //})
    //mapState({
    //user: state => state.auth.user,
    //})
    ...mapGetters("cognito", ["user"]),
  },
  methods: {
    async signUp() {
      this.loading = true;
      store
        .dispatch("cognito/signUp", {
          username: this.email,
          password: this.pass,
        })
        .then(() => {
          return store.dispatch("cognito/authenticateUser", {
            username: this.email,
            password: this.pass,
          });
        })
        .then(() => {
          this.loading = false;
          this.$router.push({ name: "home" });
        })

        .catch((e: CognitoError) => {
          this.loading = false;
          if (
            e.code == "InvalidPasswordException" ||
            e.code == "InvalidParameterException"
          ) {
            this.invalidPass = true;
            this.invalidPassMessage = e.message;
            return;
          }
          if (e.code == "UsernameExistsException") {
            this.invalidUsername = true;
            this.invalidUsernameMessage = e.message;
            return;
          }
          console.error(e);
        });
    },
  },
});
</script>

<style></style>
