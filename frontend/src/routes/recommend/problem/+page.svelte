<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import Recommend from "./Recommend.svelte";
  import type { Data } from "./data";

  export let data: Data;
</script>

<div class="flex-1 overflow-auto px-12 py-8">
  <h1 class="text-3xl">Recent Problems</h1>
  <div class="flex flex-row overflow-x-auto">
    {#each data.recent.items as item}
      <Recommend problem={item} />
    {/each}
  </div>

  {#if data.recByRating != null}
    <h1 class="mt-12 text-3xl">Recommend by Rating</h1>
    <div class="flex flex-row overflow-x-auto">
      {#each data.recByRating.items as item}
        <Recommend problem={item} />
      {/each}
    </div>
  {/if}
</div>
