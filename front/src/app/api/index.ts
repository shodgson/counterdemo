import store from "@/store";
import { RequestParams, PaymentURL, User } from "./types";

const countApi = {
  getUser(userId: string): Promise<User> {
    return sendRequest({ endpoint: `users/${userId}` });
  },
  add(userId: string, amount: number): Promise<User> {
    return sendRequest({
      endpoint: `count`,
      verb: "PATCH",
      authNeeded: true,
      body: { add: 1 },
    });
  },
  reset(userId: string): Promise<User> {
    return sendRequest({
      endpoint: `count`,
      verb: "PATCH",
      authNeeded: true,
      body: { reset: true },
    });
  },
  paymentUrl(): Promise<PaymentURL> {
    return sendRequest({
      endpoint: `subscription/activation_url`,
      authNeeded: true,
    });
  },
  cancelSubscription() {
    return sendRequest({
      endpoint: `subscription/cancel`,
      verb: "POST",
      authNeeded: true,
    });
  },
};

/*
export async function updateSeries(seriesId, series) {
  const name = await profileName();
  let endpoint = `/u/${name}/s/${seriesId}`;
  const params = await getParams(series, "PATCH", true);
  return sendRequest(endpoint, params);
}

export async function playlistInfo(playlistId) {
  let endpoint = `/playliststats/${playlistId}`;
  const params = await getParams(null, "GET", false);
  return sendRequest(endpoint, params);
}

export async function playlistSeries(playlistId) {
  let endpoint = `/playlistseries/${playlistId}`;
  const params = await getParams(null, "GET", false);
  return sendRequest(endpoint, params);
}

async function profileName() {
  console.log("Getting username")
  let profileName = store.getters["cognito/username"];
  if (profileName) {
    return profileName.replaceAll(" ", "_");
  }
  const u = await store.dispatch("cognito/getUserAttributes");
  console.debug("fetching attributes for api call");
  //console.log(u);
  return u?.preferred_username.replaceAll(" ", "_");
}
*/

async function idToken() {
  const user = await store.dispatch("cognito/getCurrentUser");
  return user.tokens?.IdToken;
}

async function parseRequest(params: RequestParams) {
  const requestHeaders: HeadersInit = new Headers();

  requestHeaders.set("Content-Type", "application/json");
  if (params.authNeeded) {
    let token = await idToken();
    //console.log(token);
    requestHeaders.set("Authorization", "Bearer " + token);
  }

  const request: RequestInit = {
    method: params.verb || "GET",
    headers: requestHeaders,
  };

  if (params.body) {
    request.body = JSON.stringify(params.body);
  }
  return request;
}

async function sendRequest(params: RequestParams) {
  const apiURL = import.meta.env.VITE_AWS_API_URL;
  const url = [apiURL, params.endpoint].join("/");
  const request = await parseRequest(params);

  return fetch(url, request).then(async (response) => {
    if (response.status == 200) {
      return response.json();
    } else {
      throw await response.text()
    }
  });
}

export default countApi;
