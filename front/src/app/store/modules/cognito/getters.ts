import { GetterTree } from "vuex";

const getters: GetterTree<any, any> = {
  isAuthenticated: (state: any) => {
    return state.user?.tokens != null;
  },
  username: (state: any) => {
    return state.user?.username;
  },
};
export default getters;
