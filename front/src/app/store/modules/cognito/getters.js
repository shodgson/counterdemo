const getters = {
  isAuthenticated: (state, getters) => {
    return state.user?.tokens != null;
  },
  username: (state, getters) => {
    return state.attributes?.preferred_username;
  },
  idToken: (state, getters) => {
    return state.user?.tokens?.IdToken;
  },
  email: (state, getters) => {
    return state.attributes?.email;
  },
};
export default getters;
