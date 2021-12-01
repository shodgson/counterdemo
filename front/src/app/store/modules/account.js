// initial state
const state = {
  active: null,
};

const getters = {
  accountActive: state => state.active,
};

const mutations = {
  activeStatus(state, status) {
    state.active = status;
  },
};

const actions = {
  setActiveStatus: (context, status) => {
    context.commit("activeStatus", status);
  },
};

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
};

