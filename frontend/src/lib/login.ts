import { readonly, writable } from "svelte/store";
import * as Event from "../lib/Messages";

export enum status {
  authorized = "authorized",
  unauthorized = "unauthorized",
  wrongpassword = "wrongpassword",
}

let store = writable(status.unauthorized);

export const auth = readonly(store);

export let conn: Event.Connection;
export function login(username: string, password: string) {
  fetch("login", {
    method: "post",
    body: JSON.stringify({
      username: username,
      password: password,
    }),
    mode: "cors",
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      } else {
        throw "unauthorized";
      }
    })
    .then((data) => {
      store.set(status.authorized);
      connectWebSocket(data.otp);
    })
    .catch((e) => {
      store.set(status.wrongpassword);
    });
}

function connectWebSocket(otp: string) {
  conn = new Event.Connection(
    new WebSocket("ws://" + document.location.host + "/ws?otp=" + otp)
  );
}
