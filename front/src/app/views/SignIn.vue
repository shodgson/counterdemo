<template>
  <div>
    <h1 class="text-xl my-6">Sign in</h1>
    <div class="w-full max-w-xs mx-auto">
      <form
        @submit.prevent="signIn"
        class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4"
      >
        <div class="mb-4">
          <label
            class="block text-gray-700 text-sm font-bold mb-2"
            for="username"
          >
            Email
          </label>
          <input
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            :class="invalidUsername ? 'border-red-500' : ''"
            @change="invalidUsername = false"
            id="emailInput"
            type="email"
            placeholder="Email"
            v-model="email"
          />
          <p v-if="invalidUsername" class="text-red-500 text-xs italic">
            {{ invalidUsernameMessage }}
          </p>
        </div>
        <div class="mb-6">
          <label
            class="block text-gray-700 text-sm font-bold mb-2"
            for="password"
          >
            Password
          </label>
          <input
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
            :class="invalidPass ? 'border-red-500' : ''"
            @change="invalidPass = false"
            id="password"
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
            Sign in
          </button>
        </div>
        <div class="mt-6">
        Need to create an account?<br/>
        <router-link :to="{ name: 'signUp'}" ><a class="link">Sign up here</a></router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import store from "@/store";
export default {
  components: {
  },
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
    async signIn() {
      this.loading = true;
      store
        .dispatch("cognito/authenticateUser", {
          username: this.email,
          password: this.pass,
        })
        .then((r) => {
          //console.log(r);
          this.loading = false;
          this.$router.push({ name: "home" });
        })
        .catch((e) => {
          this.loading = false;
          if (e.code == "InvalidPasswordException") {
            this.invalidPass = true;
            this.invalidPassMessage = e.message;
            return;
          }
          if (e.code == "NotAuthorizedException") {
            this.invalidPass = true;
            this.invalidPassMessage = e.message;
            return;
            }
          if (e.code == "UserNotFoundException") {
            this.invalidUsername = true;
            this.invalidUsernameMessage = e.message;
            return;
            }
          console.error(e);
        });
    },
  },
};
</script>

<style></style>
