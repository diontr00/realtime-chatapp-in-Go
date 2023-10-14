<script lang="ts">
  import { onMount } from "svelte";
  import { auth } from "../lib/login";
  import * as login from "../lib/login";

  let username: string;
  let password: string;
  let supported: boolean;

  onMount(() => {
    if (window["WebSocket"]) {
      supported = true;
    } else {
      supported = false;
      alert("Not supporting web socket");
    }
  });
</script>

{#if supported}
  <h1>Realtime Chat APplication</h1>
  <div style="border: 3px solid black;margin-top: 30px;">
    <form id="login-form">
      <label for="username">username:</label>
      <input
        type="text"
        id="username"
        name="username"
        bind:value={username}
      /><br />
      <label for="password">password:</label>
      <input
        type="password"
        id="password"
        name="password"
        bind:value={password}
        class={$auth === login.status.wrongpassword ? "wrongpassword" : ""}
      />

      {#if $auth === login.status.wrongpassword}
        <p style="color:red">Wrong password</p>
      {/if}
      <br /><br />
      <input
        type="submit"
        value="Login"
        class={$auth === login.status.wrongpassword ? "wrongpassword" : ""}
        on:click|preventDefault={(e) => {
          login.login(username, password);
        }}
      />
    </form>
  </div>
{:else}
  <h1>Web Socket is not supported</h1>
{/if}

<style>
  .wrongpassword {
    border: 1px solid red;
  }
</style>
