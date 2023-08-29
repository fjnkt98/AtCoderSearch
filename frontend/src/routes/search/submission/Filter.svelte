<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import Icon from "svelte-icons-pack/Icon.svelte";
  import AiOutlineSearch from "svelte-icons-pack/ai/AiOutlineSearch";
  import HiOutlinePaperAirplane from "svelte-icons-pack/hi/HiOutlinePaperAirplane";

  export let categories: string[];
  export let languages: string[];
  export let contests: string[];
  export let problems: string[];
  export let results: string[] = ["", "AC", "WA", "TLE", "RE", "CE", "MLE"];

  let category: string = "";
  let language: string = "";
  let contest: string = "";
  let problem: string = "";
  let result: string = "";
  let user: string = "";

  function searchByUser() {
    if (user !== "") {
      const params = new URLSearchParams($page.url.searchParams);
      params.set("filter.user_id", user);
      params.set("page", "1");
      goto(`${$page.url.pathname}?${params.toString()}`);
    } else {
      const params = new URLSearchParams($page.url.searchParams);
      params.delete("filter.user_id");
      params.set("page", "1");
      goto(`${$page.url.pathname}?${params.toString()}`);
    }
  }
</script>

<div class="flex flex-col items-center justify-center">
  <button
    class="mb-2 mt-4 w-full rounded-xl bg-gray-600 px-2 py-1 text-white"
    on:click={() => {
      category = "";
      language = "";
      contest = "";
      problem = "";
      result = "";
      user = "";

      const params = new URLSearchParams($page.url.searchParams);
      params.delete("filter.category");
      params.delete("filter.language");
      params.delete("filter.contest_id");
      params.delete("filter.problem_id");
      params.delete("filter.result");
      params.delete("filter.user_id");
      goto(`${$page.url.pathname}?${params.toString()}`);
    }}>Reset</button
  >

  <div class="shadow-gra-400 my-1 w-full rounded-xl bg-white px-4 py-2 shadow-md">
    <span>Contest Category</span>
    <select
      class="my-2 block w-full rounded-lg bg-gray-50 p-2.5 text-sm shadow-sm shadow-gray-300 outline-none"
      bind:value={category}
      on:change={() => {
        if (category !== "") {
          const params = new URLSearchParams($page.url.searchParams);
          params.set("filter.category", category);
          params.set("page", "1");
          goto(`${$page.url.pathname}?${params.toString()}`);
        } else {
          const params = new URLSearchParams($page.url.searchParams);
          params.delete("filter.category");
          params.set("page", "1");
          goto(`${$page.url.pathname}?${params.toString()}`);
        }
      }}
    >
      {#each categories as c}
        <option value={c}>
          {c}
        </option>
      {/each}
    </select>
  </div>

  <div class="shadow-gra-400 my-1 w-full rounded-xl bg-white px-4 py-2 shadow-md">
    <span>Language</span>
    <select
      class="my-2 block w-full rounded-lg bg-gray-50 p-2.5 text-sm shadow-sm shadow-gray-300 outline-none"
      bind:value={language}
      on:change={() => {
        if (language !== "") {
          const params = new URLSearchParams($page.url.searchParams);
          params.set("filter.language", language);
          params.set("page", "1");
          goto(`${$page.url.pathname}?${params.toString()}`);
        } else {
          const params = new URLSearchParams($page.url.searchParams);
          params.delete("filter.language");
          params.set("page", "1");
          goto(`${$page.url.pathname}?${params.toString()}`);
        }
      }}
    >
      {#each languages as l}
        <option value={l}>
          {l}
        </option>
      {/each}
    </select>
  </div>

  <div class="shadow-gra-400 my-1 w-full rounded-xl bg-white px-4 py-2 shadow-md">
    <span>Contest ID</span>
    <select
      class="my-2 block w-full rounded-lg bg-gray-50 p-2.5 text-sm shadow-sm shadow-gray-300 outline-none"
      bind:value={contest}
      on:change={() => {
        if (contest !== "") {
          const params = new URLSearchParams($page.url.searchParams);
          params.set("filter.contest_id", contest);
          params.set("page", "1");
          goto(`${$page.url.pathname}?${params.toString()}`);
        } else {
          const params = new URLSearchParams($page.url.searchParams);
          params.delete("filter.contest_id");
          params.set("page", "1");
          goto(`${$page.url.pathname}?${params.toString()}`);
        }
      }}
    >
      {#each contests as c}
        <option value={c}>
          {c}
        </option>
      {/each}
    </select>
  </div>

  <div class="shadow-gra-400 my-1 w-full rounded-xl bg-white px-4 py-2 shadow-md">
    <span>Problem ID</span>
    <select
      class="my-2 block w-full rounded-lg bg-gray-50 p-2.5 text-sm shadow-sm shadow-gray-300 outline-none"
      bind:value={problem}
      on:change={() => {
        if (problem !== "") {
          const params = new URLSearchParams($page.url.searchParams);
          params.set("filter.problem_id", problem);
          params.set("page", "1");
          goto(`${$page.url.pathname}?${params.toString()}`);
        } else {
          const params = new URLSearchParams($page.url.searchParams);
          params.delete("filter.problem_id");
          params.set("page", "1");
          goto(`${$page.url.pathname}?${params.toString()}`);
        }
      }}
    >
      {#each problems as p}
        <option value={p}>
          {p}
        </option>
      {/each}
    </select>
  </div>

  <div class="shadow-gra-400 my-1 w-full rounded-xl bg-white px-4 py-2 shadow-md">
    <span>Result</span>
    <select
      class="my-2 block w-full rounded-lg bg-gray-50 p-2.5 text-sm shadow-sm shadow-gray-300 outline-none"
      bind:value={result}
      on:change={() => {
        if (result !== "") {
          const params = new URLSearchParams($page.url.searchParams);
          params.set("filter.result", result);
          params.set("page", "1");
          goto(`${$page.url.pathname}?${params.toString()}`);
        } else {
          const params = new URLSearchParams($page.url.searchParams);
          params.delete("filter.result");
          params.set("page", "1");
          goto(`${$page.url.pathname}?${params.toString()}`);
        }
      }}
    >
      {#each results as r}
        <option value={r}>
          {r}
        </option>
      {/each}
    </select>
  </div>

  <div class="shadow-gra-400 my-1 flex w-full flex-col rounded-xl bg-white px-4 py-2 shadow-md">
    <span>User ID</span>
    <div class="white my-2 flex flex-row items-center justify-center rounded-full border border-gray-400 bg-white px-2 py-1 text-gray-900">
      <div class="px-2">
        <Icon src={AiOutlineSearch} />
      </div>
      <input
        type="text"
        class="flex-1 appearance-none bg-transparent text-gray-900 focus:border-none focus:outline-none"
        placeholder="User ID"
        on:keydown={(e) => {
          if (e.key === "Enter") {
            e.preventDefault();
            searchByUser();
          }
        }}
        bind:value={user}
      />
      <button
        type="button"
        class=""
        on:click={() => {
          searchByUser();
        }}
      >
        <Icon src={HiOutlinePaperAirplane} className="rotate-90" />
      </button>
    </div>
  </div>
</div>
