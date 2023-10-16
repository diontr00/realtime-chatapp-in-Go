import { readable, writable, readonly } from "svelte/store";
import * as message from "../model/message";
export const errors = writable("");

const closedw = writable(true);
export const closed = readonly(closedw);

let messages = writable<message.Message[]>([]);

export let incomingMessage = readonly(messages);

function appendMessage(message: message.Message) {
  messages.update((msgs) => (msgs = [...msgs, message]));
}

export function clearMessage() {
  messages.set([]);
}

export function setError(str: string) {
  errors.set(str);
}

export enum messageType {
  newMessage = "new_message",
  sendMessage = "send_message",
  errMessage = "error_message",
}

// Use to parse incoming messages
class Event {
  type: messageType;
  payload: any;
  constructor(type: messageType, payload: any) {
    this.type = type;
    this.payload = payload;
  }
}

export class Connection {
  conn: WebSocket;
  constructor(conn: WebSocket) {
    this.conn = conn;
    closedw.set(false);

    let self = this;

    conn.onmessage = function(event: MessageEvent) {
      const eventData = JSON.parse(event.data);
      const incoming = new Event(messageType.newMessage, eventData);
      self.routeEvent(incoming);
    };

    conn.onclose = function() {
      closedw.set(true);
    };
  }

  // type --  the name of the event to send
  // payload -- the data payload
  sendEvent(type: messageType, payload: message.Message): void {
    const event = new Event(type, payload);
    const byteSize = (str: string) => new Blob([str]).size;
    if (byteSize(payload.message) > 512) {
      errors.set("Message is too long");
      return;
    }
    this.conn.send(JSON.stringify(event));

    this.routeEvent(event);
  }

  // route event is a proxy method that routes
  // events into there correct handler
  private routeEvent(event: Event) {
    if (event.type === undefined) {
      alert("no 'type' field in the message");
    }
    switch (event.type) {
      case messageType.newMessage:
        appendMessage(event.payload);
        console.log("new_message", event.payload);
        break;
      case messageType.sendMessage:
        console.log("send_message", event.payload);
        break;
      case messageType.errMessage:
        errors.set(event.payload);
        console.log("error_message", event.payload);
        break;
      default:
        alert("unsupported message type");
        break;
    }
  }
}
