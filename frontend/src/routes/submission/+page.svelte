<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import Header from "$lib/Header.svelte";
  import PageNavigation from "$lib/PageNavigation.svelte";
  import Tab from "$lib/Tab.svelte";
  import type { SearchSubmissionResult } from "$lib/response";
  import { createQuery } from "@tanstack/svelte-query";
  import FilterMenu from "./FilterMenu.svelte";
  import Submission from "./Submission.svelte";
  import { fetchSearchSubmissionResult } from "./query";
  import { selections } from "./sort";

  $: width = 0;
  const widthThreshold = 768;

  let opened = false;
  let s: string = $page.url.searchParams.get("s") ?? "1";

  export let data: SearchSubmissionResult;
  const searchQuery = createQuery({
    queryKey: ["searchSubmission"],
    queryFn: async () => {
      const res = await fetchSearchSubmissionResult($page.url.searchParams, fetch);
      return res;
    },
    initialData: data,
  });
</script>

<svelte:window bind:outerWidth={width} />

<div class="flex flex-col bg-gray-100">
  <Header />
  <Tab selected={"submission"} />

  <div class={`flex  ${width > widthThreshold ? "flex-row" : "flex-col"}`}>
    <div class="flex-0 mx-2 my-1" class:max-w-96={width > widthThreshold}>
      <label class={`block cursor-pointer select-none rounded-md px-4 py-2 text-center font-semibold text-gray-50 shadow-sm shadow-gray-500 ` + (opened ? "bg-blue-800" : "bg-blue-500")}>
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
          <FilterMenu />
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
            goto(`/submission?${p.toString()}`, { replaceState: false, noScroll: true, invalidateAll: true });
          }}
        >
          {#each selections.entries() as [k, v]}
            <option value={k}>
              {v.label}
            </option>
          {/each}
        </select>

        <div class="mx-2 mb-2 flex flex-row justify-between text-sm">
          <p class="w-full text-right text-gray-700">{data.stats.count}件 約{data.stats.time}ms</p>
        </div>
      </div>

      {#if $searchQuery.isLoading}
        <p>Loading...</p>
      {:else if $searchQuery.isError}
        <p>ERROR</p>
      {:else if $searchQuery.isSuccess}
        {#each data.items as item}
          <Submission {item} />
        {/each}
        <PageNavigation current={data.stats.index} enableEnd={false} />
      {/if}
    </div>
  </div>
</div>
