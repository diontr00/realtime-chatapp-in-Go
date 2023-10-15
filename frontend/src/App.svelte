<script lang="ts">
  import * as Event from "./lib/Messages";
  import { errors, closed } from "./lib/Messages";
  import { auth } from "./lib/login";
  import * as login from "./lib/login";
  import Login from "./component/login.svelte";

  let selectedRoom = "";
  let chatroom = "";
  let outmessage = "";

  const byteSize = (str: string) => new Blob([str]).size;
  $: invalid = byteSize(outmessage) > 512 ? true : false;

  let conn: Event.Connection;

  function changeChatRoom() {
    chatroom = selectedRoom;
    console.log(chatroom);
  }

  function sendMessage() {
    if (invalid) {
      return;
    }

    if (outmessage != "") {
      conn.sendEvent(Event.messageType.sendMessage, outmessage);
    }
    outmessage = "";
  }
</script>

{#if $auth !== login.status.authorized}
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

      <h3 class={invalid ? "invalid" : "valid"} style="color: red;">
        unexpected Error: {$errors}
      </h3>
      <input
        bind:value={outmessage}
        type="text"
        id="outmessage"
        name="outmessage"
        on:keydown={() => {
          if (!invalid) {
            $errors = "";
            return;
          }
          $errors = "outmessage too long";
        }}
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
