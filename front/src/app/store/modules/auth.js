// Derived from https://github.com/wolfeidau/cognito-vue-bootstrap
// https://github.com/kopterio/vue-auth-cognito


const poolData = {
  ClientId: "4aqa7kavaig216g45l94017tt5",
  UserPoolId: "us-east-2_wvtbM3dh0", // Your user pool id here
};
import Auth from "@aws-amplify/auth";

// initial state
const state = {
  user: null,
  isAuthenticated: false,
  authenticationStatus: null
};

const getters = {
  authenticatedUser: state => state.user,
  isAuthenticated: state => state.isAuthenticated,
  authenticationStatus: state => {
    return state.authenticationStatus
      ? state.authenticationStatus
      : { variant: "secondary" };
  },
  hasAuthenticationStatus: state => {
    return !!state.authenticationStatus;
  }
};

const mutations = {
  setAuthenticationError(state, err) {
    state.authenticationStatus = {
      state: "failed",
      message: err.message,
      variant: "danger"
    };
  },
  clearAuthenticationStatus: state => {
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
  }
};

const actions = {
  clearAuthenticationStatus: context => {
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
  signOut: async context => {
    try {
      await Auth.signOut();
    } catch (err) {
    }
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
  passwordResetResend: async (context, params) => {
    context.commit("auth/clearAuthenticationStatus", null, { root: true });
    try {
      await Auth.passwordResetResend(params.username);
    } catch (err) {
      context.commit("auth/setAuthenticationError", err, { root: true });
    }
  },
  passwordChange: async (context, params) => {
    logger.debug("password change for {}", context.state.user.username);
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
  }
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
};

