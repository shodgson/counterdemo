import ActionsFactory from './actions';
import mutations from './mutations';
import getters from './getters';
import { Module } from 'vuex';

const state = {
  user: null,
  attributes: null,
};

export type CognitoConfig = {
  Region: string
UserPoolId: string;
 ClientId: string;
}

const awsconfig: CognitoConfig = {
  Region: import.meta.env.VITE_AWS_REGION as string,
  UserPoolId: import.meta.env.VITE_AWS_USER_POOL_ID as string,
  ClientId: import.meta.env.VITE_AWS_USER_POOL_CLIENT_ID as string,
};

// Example state
// const state = {
//   user: {
//     username: 'username in any format: email, UUID, etc.',
//     tokens: null | {
//       IdToken: '', // in JWT format
//       RefreshToken: '', // in JWT format
//       AccessToken: '', // in JWT format
//     },
//   },
//   attributes: {
//     email: 'user email',
//     phone_number: '+1 555 12345',
//     ...
//   },
// };

const authModule: Module<any, any> = {
  namespaced: true,
  state,
  mutations,
  getters,
}

function NewAuthModule() {
  authModule.actions = ActionsFactory(awsconfig)
  return authModule
}

export default NewAuthModule;

/*
export default class CognitoAuth {

  constructor(config: CognitoConfig) {
    this.namespaced = true;
    this.state = state;
    this.actions = new ActionsFactory(config);
    this.mutations = mutations;
    this.getters = getters;
  }
}
*/