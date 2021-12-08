<template>
  <div>
    <h1 class="text-xl my-6">Reset password</h1>
    <div class="w-full max-w-xs mx-auto">
      <form class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
        <div class="mb-4 form-control">
          <label class="label">
            Email
          </label>
          <input
            class="input input-bordered"
            :class="invalidUsername ? 'border-red-500' : ''"
            @change="invalidUsername = false"
            id="emailInput"
            type="email"
            placeholder="Email"
            v-model="email"
            :disabled="waitingForCode"
          />
          <p v-if="invalidUsername" class="text-red-500 text-xs italic">
            {{ invalidUsernameMessage }}
          </p>
        </div>
        <div v-if="!passwordReset">
          <div class="flex w-full" v-if="!waitingForCode">
            <button
              :disabled="email.length == 0"
              class="btn w-full mx-auto"
              :class="{ loading: loading }"
              @click.prevent="sendCode"
            >
              Send reset code
            </button>
          </div>
          <div class="mb-4 form-control" v-if="waitingForCode">
            <label class="label">
              <span class="label-text">Code</span>
              <span class="label-text-alt">Sent to your email</span>
            </label>
            <input
              class="input input-bordered"
              :class="invalidCode ? 'border-red-500' : ''"
              @change="invalidCode = false"
              type="text"
              placeholder="6-digit code"
              v-model="code"
            />
            <p v-if="invalidCode" class="text-red-500 text-xs italic">
              {{ invalidCodeMessage }}
            </p>
          </div>
          <div class="mb-6 form-control" v-if="waitingForCode">
            <label class="label">
              New password
            </label>
            <input
              class="input input-bordered"
              type="password"
              v-model="pass"
            />
            <p v-if="invalidPass" class="text-red-500 text-xs italic">
              {{ invalidPassMessage }}
            </p>
          </div>
          <div class="flex w-full" v-if="waitingForCode">
            <button
              :disabled="pass.length < 5 || code.length < 6"
              class="btn w-full mx-auto"
              :class="{ loading: loading }"
              @click.prevent="resetPassword"
            >
              Update password
            </button>
          </div>
        </div>
        <div v-if="passwordReset" class="alert alert-success">
          <div class="flex-col">
            <div>
            Password successfully reset
            </div>
            <div>
              <router-link :to="{ name: 'signIn' }"
                ><a class="link">Sign in</a></router-link
              >
            </div>
          </div>
        </div>
        <div class="mt-6" v-if="!passwordReset">
          <router-link :to="{ name: 'signIn' }"
            ><a class="link">Back</a></router-link
          >
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import store from "@/store";
export default {
  components: {},
  data() {
    return {
      email: "",
      pass: "",
      invalidPass: false,
      invalidPassMessage: "",
      invalidCode: false,
      invalidCodeMessage: "",
      invalidUsername: false,
      invalidUsernameMessage: "",
      loading: false,
      waitingForCode: false,
      passwordReset: false,
    };
  },
  methods: {
    async sendCode() {
      this.loading = true;
      store
        .dispatch("cognito/forgotPassword", {
          username: this.email,
        })
        .then((r) => {
          this.loading = false;
          this.waitingForCode = true;
        })
        .catch((e) => {
          this.loading = false;
          if (e.code == "UserNotFoundException") {
            this.invalidUsername = true;
            this.invalidUsernameMessage = "Email not found";
            return;
          } else {
            console.error(e);
          }
        });
    },
    resetPassword() {
      this.loading = true;
      store
        .dispatch("cognito/confirmPassword", {
          username: this.email,
          newPassword: this.pass,
          code: this.code,
        })
        .then((r) => {
          this.passwordReset = true;
        })
        .catch((e) => {
          switch (e.code) {
            case "InvalidParameterException":
              this.invalidPass = true;
              this.invalidPassMessage = "Invalid values";
              break;
            case "InvalidPasswordException":
              this.invalidPass = true;
              this.invalidPassMessage = e.message;
              break;
            case "NotAuthorizedException":
              this.invalidPass = true;
              this.invalidPassMessage = e.message;
              break;
            case "CodeMismatchException":
              this.invalidCode = true;
              this.invalidCodeMessage = e.message;
              break;
            case "LimitExceededException":
              this.invalidPass = true;
              this.invalidPassMessage = e.message;
              break;
            default:
              console.error(e);
          }
        })
        .finally(() => (this.loading = false));
    },
  },
};
</script>

<style></style>
