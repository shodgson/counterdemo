import { Module } from "vuex";

const accountModule: Module<any, any> = {
  // Initial state
  namespaced: true,

  state: {
    active: null,
  },

  getters: {
    accountActive: (state) => state.active,
  },

  mutations: {
    activeStatus(state, status: boolean) {
      state.active = status;
    },
  },

  actions: {
    setActiveStatus: (context, status) => {
      context.commit("activeStatus", status);
    },
  },
};

export default accountModule;
