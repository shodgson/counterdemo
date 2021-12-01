import * as types from './mutation-types';

export default {
  [types.AUTHENTICATE](state, payload) {
    state.user = payload;
  },
  [types.SIGNOUT](state) {
    state.user = null;
    state.attributes = null;
  },
  [types.ATTRIBUTES](state, payload) {
    state.attributes = payload;
  },
};
