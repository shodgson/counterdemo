import { ActionTree } from "vuex";
import {
  CognitoUserPool,
  CognitoUserAttribute,
  CognitoUser,
  CognitoRefreshToken,
  AuthenticationDetails,
  CognitoUserSession,
} from "amazon-cognito-identity-js";
import { CognitoConfig } from ".";

import * as types from "./mutation-types";

function constructUser(cognitoUser: CognitoUser, session: CognitoUserSession) {
  return {
    username: cognitoUser.getUsername(),
    tokens: {
      IdToken: session.getIdToken().getJwtToken(),
      AccessToken: session.getAccessToken().getJwtToken(),
      RefreshToken: session.getRefreshToken().getToken(),
    },
    attributes: {},
  };
}

// cannot use ES6 classes, the methods are not enumerable, properties are.
export default function actionsFactory(config: CognitoConfig): ActionTree<any, any> {
  const cognitoUserPool = new CognitoUserPool({
    UserPoolId: config.UserPoolId,
    ClientId: config.ClientId,
  });

  return {
    getCurrentUser({ commit }) {
      return new Promise((resolve, reject) => {
        const cognitoUser = cognitoUserPool.getCurrentUser();

        if (!cognitoUser) {
          reject({
            message: "Can't retrieve the current user",
          });
          return;
        }

        cognitoUser.getSession((err: any, session: CognitoUserSession) => {
          if (err) {
            reject(err);
            return;
          }

          const constructedUser = constructUser(cognitoUser, session);
          // Call AUTHENTICATE because it's utterly the same
          commit(types.AUTHENTICATE, constructedUser);
          resolve(constructedUser);
        });
      });
    },

    authenticateUser({ commit }, payload) {
      /* userInfo: { username, password } */
      const authDetails = new AuthenticationDetails({
        Username: payload.username,
        Password: payload.password,
      });

      const cognitoUser = new CognitoUser({
        Pool: cognitoUserPool,
        Username: payload.username,
      });

      return new Promise((resolve, reject) =>
        cognitoUser.authenticateUser(authDetails, {
          onFailure: (err) => {
            reject(err);
          },
          onSuccess: (session, userConfirmationNecessary) => {
            commit(types.AUTHENTICATE, constructUser(cognitoUser, session));
            resolve({ userConfirmationNecessary });
          },
        })
      );
    },

    signUp({ commit }, userInfo) {
      /* userInfo: { email, password, attributes } */
      const userAttributes = Object.keys(userInfo.attributes || {}).map(
        (key) =>
          new CognitoUserAttribute({
            Name: key,
            Value: userInfo.attributes[key],
          })
      );

      return new Promise((resolve, reject) => {
        cognitoUserPool.signUp(
          userInfo.username,
          userInfo.password,
          userAttributes,
          [],
          (err, data) => {
            if (!err) {
              commit(types.AUTHENTICATE, {
                username: data?.user.getUsername(),
                tokens: null, // no session yet
                attributes: {},
              });
              resolve({ userConfirmationNecessary: !data?.userConfirmed });
              return;
            }
            reject(err);
          }
        );
      });
    },

    confirmRegistration({ state }, payload) {
      /* payload: { username, code } */
      const cognitoUser = new CognitoUser({
        Pool: cognitoUserPool,
        Username: payload.username,
      });

      return new Promise<void>((resolve, reject) => {
        cognitoUser.confirmRegistration(payload.code, true, (err) => {
          if (!err) {
            resolve();
            return;
          }
          reject(err);
        });
      });
    },

    resendConfirmationCode({ commit }, payload) {
      const cognitoUser = new CognitoUser({
        Pool: cognitoUserPool,
        Username: payload.username,
      });

      return new Promise<void>((resolve, reject) => {
        cognitoUser.resendConfirmationCode((err) => {
          if (!err) {
            resolve();
            return;
          }
          reject(err);
        });
      });
    },

    forgotPassword({ commit }, payload) {
      const cognitoUser = new CognitoUser({
        Pool: cognitoUserPool,
        Username: payload.username,
      });

      return new Promise<void>((resolve, reject) =>
        cognitoUser.forgotPassword({
          onSuccess() {
            resolve();
          },
          onFailure(err) {
            reject(err);
          },
        })
      );
    },

    confirmPassword({ commit }, payload) {
      const cognitoUser = new CognitoUser({
        Pool: cognitoUserPool,
        Username: payload.username,
      });

      return new Promise<void>((resolve, reject) => {
        cognitoUser.confirmPassword(payload.code, payload.newPassword, {
          onFailure(err) {
            reject(err);
          },
          onSuccess() {
            resolve();
          },
        });
      });
    },

    // Only for authenticated users
    changePassword({ state }, payload) {
      return new Promise<void>((resolve, reject) => {
        // Make sure the user is authenticated
        if (state.user === null || (state.user && state.user.tokens === null)) {
          reject({
            message: "User is unauthenticated",
          });
          return;
        }

        const cognitoUser = new CognitoUser({
          Pool: cognitoUserPool,
          Username: state.user.username,
        });

        // Restore session without making an additional call to API
        //cognitoUser.signInUserSession = cognitoUser.getCognitoUserSession(
        //  state.user.tokens
        //);

        cognitoUser.changePassword(
          payload.oldPassword,
          payload.newPassword,
          (err) => {
            if (!err) {
              resolve();
              return;
            }
            reject(err);
          }
        );
      });
    },

    // Only for authenticated users
    updateAttributes({ commit, state }, payload) {
      /* payload: { [attributes] } */
      return new Promise<void>((resolve, reject) => {
        // Make sure the user is authenticated
        if (state.user === null || (state.user && state.user.tokens === null)) {
          reject({
            message: "User is unauthenticated",
          });
          return;
        }

        const cognitoUser = new CognitoUser({
          Pool: cognitoUserPool,
          Username: state.user.username,
        });

        // Restore session without making an additional call to API
        //cognitoUser.signInUserSession = cognitoUser.getCognitoUserSession(
          //state.user.tokens
        //);

        const attributes = Object.keys(payload || {}).map(
          (key) =>
            new CognitoUserAttribute({
              Name: key,
              Value: payload[key],
            })
        );

        cognitoUser.updateAttributes(attributes, (err) => {
          if (!err) {
            resolve();
            return;
          }
          reject(err);
        });
      });
    },

    // Only for authenticated users
    getUserAttributes({ commit, state }) {
      return new Promise((resolve, reject) => {
        // Make sure the user is authenticated
        if (state.user === null || (state.user && state.user.tokens === null)) {
          reject({
            message: "User is unauthenticated",
          });
          return;
        }

        const cognitoUser = new CognitoUser({
          Pool: cognitoUserPool,
          Username: state.user.username,
        });

        // Restore session without making an additional call to API
        //cognitoUser.signInUserSession = cognitoUser.getCognitoUserSession(
          //state.user.tokens
        //);

        cognitoUser.getUserAttributes((err, attributes) => {
          if (err) {
            reject(err);
            return;
          }

          attributes = attributes || []
          const attributesMap = new Map(attributes.map(i => [i.Name, i.Value]))
          

          //const attributesMap = (attributes || []).reduce((accum, item) => {
            //accum[item.Name] = item.Value;
            //return accum;
          //}, {});

          commit(types.ATTRIBUTES, attributesMap);
          resolve(attributesMap);
        });
      });
    },

    // Only for authenticated users
    signOut({ commit, state }) {
      return new Promise<void>((resolve, reject) => {
        // Make sure the user is authenticated
        if (state.user === null || (state.user && state.user.tokens === null)) {
          reject({
            message: "User is unauthenticated",
          });
          return;
        }

        const cognitoUser = new CognitoUser({
          Pool: cognitoUserPool,
          Username: state.user.username,
        });

        cognitoUser.signOut();
        commit(types.SIGNOUT);
        resolve();
      });
    },

    // Only for authenticated users
    refreshSession({ state, commit }) {
      return new Promise<void>((resolve, reject) => {
        // Make sure the user is authenticated
        const cognitoUser = cognitoUserPool.getCurrentUser();
        if (state.user === null || (state.user && state.user.tokens === null)) {
          reject({
            message: "User is unauthenticated",
          })
        }

        let token = new CognitoRefreshToken({
          RefreshToken: state.user.tokens.RefreshToken,
        })

        cognitoUser?.refreshSession(token, (err, session) => {
          if (!err) {
            commit(types.AUTHENTICATE, constructUser(cognitoUser, session));
            resolve();
          }
          reject(err);
        })
      })
    },
  };
}

