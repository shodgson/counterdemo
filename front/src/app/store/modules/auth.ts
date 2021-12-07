import { Module } from "vuex";
import Auth from "@aws-amplify/auth";

// Derived from https://github.com/wolfeidau/cognito-vue-bootstrap
// https://github.com/kopterio/vue-auth-cognito

const authModule: Module<any, any> = {
  // initial state
  state: {
    user: null,
    isAuthenticated: false,
    authenticationStatus: null,
  },

  getters: {
    authenticatedUser: (state) => state.user,
    isAuthenticated: (state) => state.isAuthenticated,
    authenticationStatus: (state) => {
      return state.authenticationStatus
        ? state.authenticationStatus
        : { variant: "secondary" };
    },
    hasAuthenticationStatus: (state) => {
      return !!state.authenticationStatus;
    },
  },

  mutations: {
    setAuthenticationError(state, err) {
      state.authenticationStatus = {
        state: "failed",
        message: err.message,
        variant: "danger",
      };
    },
    clearAuthenticationStatus: (state) => {
      state.authenticationStatus = null;
    },
    setUserAuthenticated(state, user) {
      state.user = user;
      state.isAuthenticated = true;
    },
    clearAuthentication(state) {
      state.user = null;
      state.userId = null;
      state.isAuthenticated = false;
    },
  },

  actions: {
    clearAuthenticationStatus: (context) => {
      context.commit("clearAuthenticationStatus", null);
    },
    signIn: async (context, params) => {
      context.commit("auth/clearAuthenticationStatus", null, { root: true });
      try {
        const user = await Auth.signIn(params.username, params.password);
        context.commit("setUserAuthenticated", user);
      } catch (err) {
        context.commit("auth/setAuthenticationError", err, { root: true });
      }
    },
    signOut: async (context) => {
      try {
        await Auth.signOut();
      } catch (err) {}
      context.commit("auth/clearAuthentication", null, { root: true });
    },
    signUp: async (context, params) => {
      context.commit("auth/clearAuthenticationStatus", null, { root: true });
      try {
        await Auth.signUp(params);
        context.commit("auth/clearAuthentication", null, { root: true });
      } catch (err) {
        context.commit("auth/setAuthenticationError", err, { root: true });
      }
    },
    confirmSignUp: async (context, params) => {
      context.commit("auth/clearAuthenticationStatus", null, { root: true });
      try {
        await Auth.confirmSignUp(params.username, params.code);
      } catch (err) {
        context.commit("auth/setAuthenticationError", err, { root: true });
      }
    },
    confirmResend: async (context, params) => {
      context.commit("auth/clearAuthenticationStatus", null, { root: true });
      try {
        await Auth.resendSignUp(params.username);
      } catch (err) {
        context.commit("auth/setAuthenticationError", err, { root: true });
      }
    },
    passwordReset: async (context, params) => {
      context.commit("auth/clearAuthenticationStatus", null, { root: true });
      try {
        await Auth.forgotPassword(params.username);
      } catch (err) {
        context.commit("auth/setAuthenticationError", err, { root: true });
      }
    },
    confirmPasswordReset: async (context, params) => {
      context.commit("auth/clearAuthenticationStatus", null, { root: true });
      try {
        await Auth.forgotPasswordSubmit(
          params.username,
          params.code,
          params.password
        );
      } catch (err) {
        context.commit("auth/setAuthenticationError", err, { root: true });
      }
    },
    passwordChange: async (context, params) => {
      context.commit("auth/clearAuthenticationStatus", null, { root: true });
      try {
        const user = await Auth.currentAuthenticatedUser();
        await Auth.changePassword(
          user,
          params.currentPassword,
          params.newPassword
        );
      } catch (err) {
        context.commit("auth/setAuthenticationError", err, { root: true });
      }
    },
  },
};
export default authModule;
