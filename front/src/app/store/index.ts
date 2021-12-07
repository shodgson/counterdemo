import { createStore } from "vuex";

// Modules
import CognitoAuth from "./modules/cognito";
import account from "./modules/account";

const store = createStore({
  modules: {
    cognito: CognitoAuth(),
    account,
  },
});

export default store;
