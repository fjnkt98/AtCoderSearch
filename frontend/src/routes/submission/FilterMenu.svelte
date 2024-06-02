<script lang="ts">
  import { browser } from "$app/environment";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import { env } from "$env/dynamic/public";
  import type { ListLanguageResponse, ListResponse } from "$lib/response";
  import { createQuery } from "@tanstack/svelte-query";

  let category: string | null = $page.url.searchParams.get("category");
  let language: string | null = $page.url.searchParams.get("language");
  let contestId: string | null = $page.url.searchParams.get("contestId");
  let problemId: string | null = $page.url.searchParams.get("problemId");
  let userId: string | null = $page.url.searchParams.get("userId");
  let result: string | null = $page.url.searchParams.get("result");
  let pointFrom: string | null = $page.url.searchParams.get("pointFrom");
  let pointTo: string | null = $page.url.searchParams.get("pointTo");
  let lengthFrom: string | null = $page.url.searchParams.get("lengthFrom");
  let lengthTo: string | null = $page.url.searchParams.get("lengthTo");
  let executionTimeFrom: string | null = $page.url.searchParams.get("executionTimeFrom");
  let executionTimeTo: string | null = $page.url.searchParams.get("executionTimeTo");

  function filter() {
    const p = new URLSearchParams($page.url.searchParams);

    if (category != null) {
      p.set("category", category);
    } else {
      p.delete("category");
    }
    if (language != null) {
      p.set("languageGroup", language);
    } else {
      p.delete("languageGroup");
    }
    if (contestId != null) {
      p.set("contestId", contestId);
    } else {
      p.delete("contestId");
    }
    if (problemId != null) {
      p.set("problemId", problemId);
    } else {
      p.delete("problemId");
    }
    if (userId != null && userId !== "") {
      p.set("userId", userId);
    } else {
      p.delete("userId");
    }
    if (result != null) {
      p.set("result", result);
    } else {
      p.delete("result");
    }
    if (pointFrom != null && pointFrom !== "") {
      p.set("pointFrom", pointFrom);
    } else {
      p.delete("pointFrom");
    }
    if (pointTo != null && pointTo !== "") {
      p.set("pointTo", pointTo);
    } else {
      p.delete("pointTo");
    }
    if (lengthFrom != null && lengthFrom !== "") {
      p.set("lengthFrom", lengthFrom);
    } else {
      p.delete("lengthFrom");
    }
    if (lengthTo != null && lengthTo !== "") {
      p.set("lengthTo", lengthTo);
    } else {
      p.delete("lengthTo");
    }
    if (executionTimeFrom != null && executionTimeFrom !== "") {
      p.set("executionTimeFrom", executionTimeFrom);
    } else {
      p.delete("executionTimeFrom");
    }
    if (executionTimeTo != null && executionTimeTo !== "") {
      p.set("executionTimeTo", executionTimeTo);
    } else {
      p.delete("executionTimeTo");
    }

    goto(`/submission?${p.toString()}`, { replaceState: false, noScroll: true, invalidateAll: true });
  }

  const categoryQuery = createQuery({
    queryKey: ["listCategory"],
    queryFn: async () => {
      const host = browser ? String(env.PUBLIC_EXTERNAL_API_HOST) : String(env.PUBLIC_INTERNAL_API_HOST);
      const response = await fetch(`${host}/api/list/category`);
      const result: ListResponse = await response.json();
      return result;
    },
  });
  const languageQuery = createQuery({
    queryKey: ["listLanguage"],
    queryFn: async () => {
      const host = browser ? String(env.PUBLIC_EXTERNAL_API_HOST) : String(env.PUBLIC_INTERNAL_API_HOST);
      const response = await fetch(`${host}/api/list/language`);
      const result: ListLanguageResponse = await response.json();
      return result;
    },
  });
  $: contestQuery = createQuery({
    queryKey: ["listContest", category],
    queryFn: async () => {
      const p = new URLSearchParams();
      if (category != null) {
        p.set("category", category);
      } else {
        p.delete("category");
      }
      const host = browser ? String(env.PUBLIC_EXTERNAL_API_HOST) : String(env.PUBLIC_INTERNAL_API_HOST);
      const response = await fetch(`${host}/api/list/contest?${p.toString()}`);
      const result: ListResponse = await response.json();
      return result;
    },
  });
  $: problemQuery = createQuery({
    queryKey: ["listProblem", contestId],
    queryFn: async () => {
      const p = new URLSearchParams();
      if (contestId != null) {
        p.set("contestId", contestId);
      } else {
        p.delete("contestId");
      }
      if (category != null) {
        p.set("category", category);
      } else {
        p.delete("category");
      }
      const host = browser ? String(env.PUBLIC_EXTERNAL_API_HOST) : String(env.PUBLIC_INTERNAL_API_HOST);
      const response = await fetch(`${host}/api/list/problem?${p.toString()}`);
      const result: ListResponse = await response.json();
      return result;
    },
  });
</script>

<div class="flex flex-wrap gap-2 bg-white px-2 py-2 text-sm">
  {#if $categoryQuery.isSuccess}
    <div class="my-1 rounded-md px-2 py-1 shadow-sm shadow-gray-500">
      <p class="text-md font-semibold">カテゴリ</p>
      <select
        class="block rounded-md border border-gray-400 bg-white px-2 py-1"
        bind:value={category}
        on:change={() => {
          filter();
        }}
      >
        <option value={null}> </option>
        {#each $categoryQuery.data.items as c}
          <option value={c}>
            {c}
          </option>
        {/each}
      </select>
    </div>
  {/if}

  {#if $languageQuery.isSuccess}
    <div class="my-1 rounded-md px-2 py-1 shadow-sm shadow-gray-500">
      <p class="text-md font-semibold">言語</p>
      <select
        class="block rounded-md border border-gray-400 bg-white px-2 py-1"
        bind:value={language}
        on:change={() => {
          filter();
        }}
      >
        <option value={null}> </option>
        {#each $languageQuery.data.items as l}
          <option value={l.group}>
            {l.group}
          </option>
        {/each}
      </select>
    </div>
  {/if}

  {#if $contestQuery.isSuccess}
    <div class="my-1 rounded-md px-2 py-1 shadow-sm shadow-gray-500">
      <p class="text-md font-semibold">コンテスト</p>
      <select
        class="block rounded-md border border-gray-400 bg-white px-2 py-1"
        bind:value={contestId}
        on:change={() => {
          filter();
        }}
      >
        <option value={null}> </option>
        {#each $contestQuery.data.items as c}
          <option value={c}>
            {c}
          </option>
        {/each}
      </select>
    </div>
  {/if}

  {#if $problemQuery.isSuccess}
    <div class="my-1 rounded-md px-2 py-1 shadow-sm shadow-gray-500">
      <p class="text-md font-semibold">問題</p>
      <select
        class="block rounded-md border border-gray-400 bg-white px-2 py-1"
        bind:value={contestId}
        on:change={() => {
          filter();
        }}
      >
        <option value={null}> </option>
        {#each $problemQuery.data.items as p}
          <option value={p}>
            {p}
          </option>
        {/each}
      </select>
    </div>
  {/if}

  <div class="my-1 rounded-md px-2 py-1 shadow-sm shadow-gray-500">
    <p class="text-md font-semibold">結果</p>
    <select
      class="block rounded-md border border-gray-400 bg-white px-2 py-1"
      bind:value={result}
      on:change={() => {
        filter();
      }}
    >
      <option value={null}> </option>
      <option value="AC">AC</option>
      <option value="WA">WA</option>
      <option value="TLE">TLE</option>
      <option value="RE">RE</option>
      <option value="CE">CE</option>
      <option value="MLE">MLE</option>
    </select>
  </div>

  <div class="my-1 rounded-md px-2 py-1 shadow-sm shadow-gray-500">
    <p class="text-md font-semibold">ユーザID</p>
    <input
      type="search"
      class="rounded-md border border-gray-900 bg-transparent px-2 text-gray-700 focus:border-blue-700"
      placeholder="ユーザID"
      bind:value={userId}
      on:keydown={(e) => {
        if (e.key === "Enter") {
          e.preventDefault();
          filter();
        }
      }}
    />
  </div>

  <div class="my-1 rounded-md px-2 py-1 shadow-sm shadow-gray-500">
    <p class="text-md font-semibold">得点</p>
    <input
      type="number"
      class="max-w-24"
      bind:value={pointFrom}
      on:keydown={(e) => {
        if (e.key === "Enter") {
          e.preventDefault();
          filter();
        }
      }}
    />
    <span class="mx-1">〜</span>
    <input
      type="number"
      class="max-w-24"
      bind:value={pointTo}
      on:keydown={(e) => {
        if (e.key === "Enter") {
          e.preventDefault();
          filter();
        }
      }}
    />
  </div>

  <div class="my-1 rounded-md px-2 py-1 shadow-sm shadow-gray-500">
    <p class="text-md font-semibold">コード長</p>
    <input
      type="number"
      class="max-w-24"
      bind:value={lengthFrom}
      on:keydown={(e) => {
        if (e.key === "Enter") {
          e.preventDefault();
          filter();
        }
      }}
    />
    <span class="mx-1">〜</span>
    <input
      type="number"
      class="max-w-24"
      bind:value={lengthTo}
      on:keydown={(e) => {
        if (e.key === "Enter") {
          e.preventDefault();
          filter();
        }
      }}
    />
  </div>

  <div class="my-1 rounded-md px-2 py-1 shadow-sm shadow-gray-500">
    <p class="text-md font-semibold">実行時間</p>
    <input
      type="number"
      class="max-w-24"
      bind:value={executionTimeFrom}
      on:keydown={(e) => {
        if (e.key === "Enter") {
          e.preventDefault();
          filter();
        }
      }}
    />
    <span class="mx-1">〜</span>
    <input
      type="number"
      class="max-w-24"
      bind:value={executionTimeTo}
      on:keydown={(e) => {
        if (e.key === "Enter") {
          e.preventDefault();
          filter();
        }
      }}
    />
  </div>
</div>
