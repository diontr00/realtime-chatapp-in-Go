import { readonly, writable } from "svelte/store";
import * as Event from "../lib/Messages";

export enum status {
  authorized = "authorized",
  unauthorized = "unauthorized",
  wrongpassword = "wrongpassword",
  badrequest = "badrequest",
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
    headers: { "Content-Type": "application/json" },
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      } else {
        console.log(response.json());
        if (response.status === 400) {
          throw "bad request";
        }
        throw "unauthorized";
      }
    })
    .then((data) => {
      connectWebSocket(data.otp.Key);

      store.set(status.authorized);
    })
    .catch((e) => {
      if (e === "bad request") {
        store.set(status.badrequest);
        return;
      }
      store.set(status.wrongpassword);
    });
}

function connectWebSocket(otp: string) {
  conn = new Event.Connection(
    new WebSocket("wss://" + document.location.host + "/ws?otp=" + otp)
  );
}
