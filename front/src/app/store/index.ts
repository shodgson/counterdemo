import { createStore } from "vuex";

// Modules
import CognitoAuth from "./modules/cognito";
import cognitoConfig from "@/config/cognito.ts";
import account from "./modules/account";

const store = new createStore({
  modules: {
    cognito: new CognitoAuth(cognitoConfig),
    account,
  },
});

export default store;
