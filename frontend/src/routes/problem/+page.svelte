<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import Header from "$lib/Header.svelte";
  import PageNavigation from "$lib/PageNavigation.svelte";
  import SearchBar from "$lib/SearchBar.svelte";
  import Tab from "$lib/Tab.svelte";
  import type { SearchProblemResult } from "$lib/response";
  import { createQuery } from "@tanstack/svelte-query";
  import FilterMenu from "./FilterMenu.svelte";
  import Problem from "./Problem.svelte";
  import { fetchSearchProblemResult } from "./query";
  import { selections } from "./sort";

  $: width = 0;
  const widthThreshold = 768;

  let opened = false;
  let s: string = $page.url.searchParams.get("s") ?? "2";

  export let data: SearchProblemResult;
  const searchQuery = createQuery({
    queryKey: ["searchProblem"],
    queryFn: async () => {
      const res = await fetchSearchProblemResult($page.url.searchParams, fetch);
      return res;
    },
    initialData: data,
  });
</script>

<svelte:window bind:outerWidth={width} />

<div class="flex flex-col bg-gray-100">
  <Header />
  <Tab selected={"problem"} />

  <div class={`flex  ${width > widthThreshold ? "flex-row" : "flex-col"}`}>
    <div class="flex-0 mx-2 my-1 min-w-32">
      <label class={`block cursor-pointer select-none rounded-md py-2 text-center font-semibold text-gray-50 shadow-sm shadow-gray-500 ` + (opened ? "bg-blue-800" : "bg-blue-500")}>
        絞り込み
        <input
          type="button"
          class="hidden"
          on:click={() => {
            opened = !opened;
          }}
        />
      </label>
      {#if $searchQuery.isLoading}
        <p>Loading...</p>
      {:else if $searchQuery.isError}
        <p>ERROR</p>
      {:else if $searchQuery.isSuccess}
        <div class:hidden={!opened}>
          <FilterMenu facet={data.stats.facet} />
        </div>
      {/if}
    </div>

    <div class="min-w-0 flex-grow">
      <div class="my-1 flex flex-row items-center justify-between gap-2 px-2 text-sm">
        <select
          class="block rounded-md border border-gray-400 bg-white px-2 py-1"
          bind:value={s}
          on:change={() => {
            const p = $page.url.searchParams;
            p.set("s", s);
            goto(`/problem?${p.toString()}`, { replaceState: false, noScroll: true, invalidateAll: true });
          }}
        >
          {#each selections.entries() as [k, v]}
            <option value={k}>
              {v.label}
            </option>
          {/each}
        </select>

        <div class="w-2/3">
          <SearchBar
            q={$page.url.searchParams.get("q") ?? ""}
            on:search={(e) => {
              const p = new URLSearchParams();
              if (e.detail == null) {
                p.delete("q");
              } else {
                p.set("q", e.detail);
              }
              const s = $page.url.searchParams.get("s");
              if (s != null) {
                p.set("s", s);
              }
              goto(`/problem?${p.toString()}`, { replaceState: false, keepFocus: true, invalidateAll: true });
            }}
          />
        </div>
      </div>

      <div class="mx-2 mb-2 flex flex-row justify-between text-sm">
        <p class="w-full text-right text-gray-700">{data.stats.count}件/{data.stats.total}件 約{data.stats.time}ms</p>
      </div>

      {#if $searchQuery.isLoading}
        <p>Loading...</p>
      {:else if $searchQuery.isError}
        <p>ERROR</p>
      {:else if $searchQuery.isSuccess}
        {#each data.items as item}
          <Problem {item} />
        {/each}
        <PageNavigation current={data.stats.index} end={data.stats.pages} />
      {/if}
    </div>
  </div>
</div>
