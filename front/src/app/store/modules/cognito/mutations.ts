import * as types from './mutation-types';

export default {
  [types.AUTHENTICATE](state: any, payload: any) {
    state.user = payload;
  },
  [types.SIGNOUT](state: any) {
    state.user = null;
    state.attributes = null;
  },
  [types.ATTRIBUTES](state: any, payload: any) {
    state.attributes = payload;
  },
};
