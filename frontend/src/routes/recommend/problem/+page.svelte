<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import Recommend from "./Recommend.svelte";
  import type { Data } from "./data";

  export let data: Data;

  let difficulty = "normal";
  let category = "ABC";
  let trend = "recent";
  let nonexperimental = false;
  let unsolved = false;

  function search() {
    const params = new URLSearchParams($page.url.searchParams);
    const option = [...(params.get("option") ?? "1111")];

    if (difficulty === "easy") {
      option[0] = "0";
    } else if (difficulty === "normal") {
      option[0] = "1";
    } else if (difficulty === "hard") {
      option[0] = "2";
    }

    if (category === "ALL") {
      option[1] = "0";
    } else if (category === "ABC") {
      option[1] = "1";
    } else if (category === "ARC") {
      option[1] = "2";
    } else if (category === "AGC") {
      option[1] = "3";
    }

    if (trend === "all") {
      option[2] = "0";
    } else if (trend === "recent") {
      option[2] = "1";
    }

    if (nonexperimental) {
      option[3] = "1";
    } else {
      option[3] = "0";
    }

    params.set("option", option.join(""));

    goto(`${$page.url.pathname}?${params.toString()}`, { replaceState: false });
  }
</script>

<div class="flex-1 overflow-auto px-4 py-3 sm:px-12 sm:py-8">
  <h1 class="text-2xl sm:text-3xl">Recent Problems</h1>
  <div class="flex flex-row overflow-x-auto">
    {#each data.recent.items as item}
      <Recommend problem={item} />
    {/each}
  </div>

  {#if data.recByRating != null}
    <h1 class="mb-2 mt-6 text-2xl sm:mb-4 sm:mt-12 sm:text-3xl">Recommend by Rating</h1>

    <div class="flex flex-row flex-wrap items-center">
      <div class="m-1 rounded-xl bg-white px-4 py-0.5 shadow-sm shadow-gray-400 sm:py-2">
        <label class="mx-1 cursor-pointer select-none">
          <input type="radio" bind:group={difficulty} value="easy" on:change={() => search()} />
          easy
        </label>
        <label class="mx-1 cursor-pointer select-none">
          <input type="radio" bind:group={difficulty} value="normal" on:change={() => search()} />
          normal
        </label>
        <label class="mx-1 cursor-pointer select-none">
          <input type="radio" bind:group={difficulty} value="hard" on:change={() => search()} />
          hard
        </label>
      </div>

      <div class="m-1 rounded-xl bg-white px-4 py-0.5 shadow-sm shadow-gray-400 sm:py-2">
        <label class="mx-1 cursor-pointer select-none">
          <input type="radio" bind:group={category} value="ALL" on:change={() => search()} />
          ALL
        </label>
        <label class="mx-1 cursor-pointer select-none">
          <input type="radio" bind:group={category} value="ABC" on:change={() => search()} />
          ABC
        </label>
        <label class="mx-1 cursor-pointer select-none">
          <input type="radio" bind:group={category} value="ARC" on:change={() => search()} />
          ARC
        </label>
        <label class="mx-1 cursor-pointer select-none">
          <input type="radio" bind:group={category} value="AGC" on:change={() => search()} />
          AGC
        </label>
      </div>

      <div class="m-1 rounded-xl bg-white px-4 py-0.5 shadow-sm shadow-gray-400 sm:py-2">
        <label class="mx-1 cursor-pointer select-none">
          <input type="radio" bind:group={trend} value="recent" on:change={() => search()} />
          recent
        </label>
        <label class="mx-1 cursor-pointer select-none">
          <input type="radio" bind:group={trend} value="all" on:change={() => search()} />
          all
        </label>
      </div>

      <div class="m-1 rounded-xl bg-white px-4 py-0.5 shadow-sm shadow-gray-400 sm:py-2">
        <label class="mx-1 cursor-pointer select-none">
          <input type="checkbox" bind:checked={nonexperimental} on:change={() => search()} />
          prefer not-experimental
        </label>
        <label class="mx-1 cursor-pointer select-none">
          <input type="checkbox" bind:checked={unsolved} on:change={() => search()} />
          exclude solved
        </label>
      </div>
    </div>

    <div class="flex flex-row overflow-x-auto">
      {#each data.recByRating.items as item}
        <Recommend problem={item} />
      {/each}
    </div>
  {/if}
</div>
