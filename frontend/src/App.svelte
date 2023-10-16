<script lang="ts">
  import * as Event from "./lib/Messages";
  import {
    errors,
    closed,
    incomingMessage,
    clearMessage,
  } from "./lib/Messages";
  import { auth } from "./lib/login";
  import * as login from "./lib/login";
  import Login from "./component/login.svelte";
  import { NewMessage } from "./model/message";
  import { conn } from "./lib/login";

  let selectedRoom = "";
  let chatroom = "";
  let outmessage = "";
  let inmessage = "";
  let textArea: HTMLElement;

  let invalid: boolean = false;

  function checkByteSize() {
    invalid = byteSize(outmessage) > 512 ? true : false;

    if (invalid && $errors === "") {
      $errors = "outmessage too long";
    }
    if (!invalid && $errors === "outmessage too long") {
      $errors = "";
    }
  }

  function checkRoom(): boolean {
    invalid = chatroom === "" ? true : false;

    if (invalid && $errors === "") {
      $errors = "room missing";
      return false;
    }

    if (!invalid) {
      return true;
    }
  }

  const byteSize = (str: string) => new Blob([str]).size;

  function changeChatRoom() {
    if ($errors === "room missing") {
      $errors = "";
    }
    chatroom = selectedRoom;
  }

  function sendMessage() {
    if (invalid) {
      return;
    }
    if (!checkRoom()) {
      return;
    }

    if (outmessage != "") {
      let message = NewMessage(chatroom, $auth.user, outmessage);

      conn.sendEvent(Event.messageType.sendMessage, message);

      outmessage = "";
    }
  }

  $: if ($incomingMessage.length > 0) {
    textArea = document.getElementById("chatmessages");

    if (textArea !== null) {
      for (let message of $incomingMessage) {
        if (message.room === chatroom) {
          const formattedMsg = `${message.user}: ${message.message}`;
          inmessage += formattedMsg + "\n";
        }
      }
      textArea.innerHTML = textArea.innerHTML + "\n" + inmessage;
      clearMessage();
    }
  }
</script>

{#if $errors}
  <h3 style="color:red">{$errors}</h3>
{/if}

{#if $auth.status !== login.status.authorized}
  <Login />
{:else}
  {#if $closed}
    <h1 style="color:red">Connection have been closed</h1>
  {/if}
  <main>
    <h1>Chat Application</h1>

    {#if chatroom !== ""}
      <h3 id="chat-header">Currently in chat: {chatroom}</h3>
    {/if}

    <label for="Chatroom">Choose a chat room:</label>

    <select name="chatroom" bind:value={selectedRoom} required>
      <option value="general">general</option>
      <option value="special1">special1</option>
    </select>

    <button on:click|preventDefault={changeChatRoom}>Change Chatroom</button>

    <br />

    {#if chatroom !== ""}
      <textarea
        bind:this={textArea}
        class="messagearea"
        name="chatmessages"
        id="chatmessages"
        readonly
        cols="50"
        rows="4"
        placeholder="Welcome to the {chatroom} chatroom, messages from other will appear here"
      />
    {:else}
      <h2>Change chatroom to receive outmessage</h2>
    {/if}
    <h3>Input Message:</h3>

    <!-- Form use to send outmessage -->
    <form id="chatroom-message">
      <label for="outmessage">Message:</label>

      <input
        bind:value={outmessage}
        type="text"
        id="outmessage"
        name="outmessage"
        on:keydown={checkByteSize}
        class={invalid ? "invalid" : ""}
      />
      <br /><br />
      <button
        on:click|preventDefault={sendMessage}
        class={invalid ? "invalid" : ""}>Send Message</button
      >
    </form>
  </main>
{/if}

<style>
  .invalid {
    border: 1px solid red;
  }

  .valid {
    display: none;
  }
</style>
