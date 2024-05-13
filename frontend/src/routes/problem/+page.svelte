<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import Header from "$lib/Header.svelte";
  import SearchBar from "$lib/SearchBar.svelte";
  import Tab from "$lib/Tab.svelte";
  import { bgColorStyles, textColorStyles } from "$lib/colors";
  import type { SearchProblemResult } from "$lib/response";
  import { createQuery } from "@tanstack/svelte-query";
  import dayjs from "dayjs";
  import timezone from "dayjs/plugin/timezone";
  import utc from "dayjs/plugin/utc";
  import FilterMenu from "./FilterMenu.svelte";
  import { fetchSearchProblemResult, fetchRecommendProblemResult } from "./query";
  import { selections } from "./sort";

  dayjs.extend(timezone);
  dayjs.extend(utc);

  $: width = 0;
  const widthThreshold = 768;

  let opened = false;
  $: {
    if (width > widthThreshold) {
      opened = true;
    } else {
      opened = false;
    }
  }
  let detail: string | null = "abc351_a";
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

  $: recommendQuery = createQuery({
    queryKey: ["recommendProblem", detail],
    queryFn: async () => {
      const res = await fetchRecommendProblemResult({ limit: 5, page: 1, problemId: detail }, fetch);
      return res;
    },
  });

  function convertDateTime(date: string): string {
    return dayjs(date).tz("Asia/Tokyo").format("YYYY/MM/DD HH:mm:ss");
  }
</script>

<svelte:window bind:outerWidth={width} />

<div class="flex h-dvh flex-col bg-gray-100">
  <Header />
  <Tab selected={"problem"} />

  <div class={`flex  ${width > widthThreshold ? "flex-row" : "flex-col"}`}>
    <div class="flex-0 mx-2 my-1">
      <label class={`block select-none rounded-md py-2 text-center font-semibold shadow-sm shadow-gray-500 ` + (opened ? "bg-blue-500 text-gray-50" : "text-gray-800")}>
        絞り込み
        <input
          type="button"
          class="hidden"
          on:click={() => {
            if (width > widthThreshold) {
              opened = true;
            } else {
              opened = !opened;
            }
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

    <div class="flex-grow border border-red-500">
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
          <div class="mx-2 my-0.5 flex flex-col rounded-lg bg-white px-3 py-3 shadow-md shadow-gray-200">
            <div class="flex flex-row items-center justify-between">
              <a href={item.contestUrl} target="_blank" class="block font-medium text-blue-500">
                {item.contestId.toUpperCase()}
              </a>
              <a href={item.problemUrl} target="_blank" class={`block text-pretty font-medium ${textColorStyles.get(item.color ?? "black")}`}>
                {item.problemTitle}
              </a>
              <label class={`w-12 text-sm text-black ` + (detail == item.problemId ? "" : "underline")}>
                {detail == item.problemId ? "閉じる" : "詳細"}
                <input
                  type="button"
                  class="h-0 w-0"
                  on:click={() => {
                    detail = detail == item.problemId ? null : item.problemId;
                  }}
                />
              </label>
            </div>
            {#if detail == item.problemId}
              <div class="flex flex-row py-1">
                <div class="my-0.5 flex flex-grow flex-col items-start justify-between">
                  <span class="block text-sm text-slate-500">{item.contestTitle}</span>
                  <span class="block text-sm text-slate-500">{convertDateTime(item.startAt)}</span>
                </div>
                <div class="flex flex-col items-center justify-between">
                  <span class={`block rounded-full px-2 py-0.5 text-xs font-medium text-white  ${bgColorStyles.get(item.category)}`}>{item.category}</span>
                  {#if item.difficulty != null}
                    <span class={`block rounded-full px-2 py-0.5 text-xs font-medium text-white ${bgColorStyles.get(item.color ?? "black")}`}>{item.difficulty}</span>
                  {/if}
                </div>
              </div>

              {#if $recommendQuery.isSuccess}
                <p class="">似ているかも</p>
                <div class="my-1 flex flex-row gap-4 overflow-x-scroll py-0.5">
                  {#each $recommendQuery.data.items as item}
                    <div class="flex flex-none flex-col text-nowrap rounded-md px-3 py-2 shadow-md shadow-gray-500">
                      <p class="block text-nowrap text-sm text-gray-700">{item.contestId.toUpperCase()}</p>
                      <a href={item.problemUrl} class="text-md block flex-grow text-wrap text-gray-700" target="_blank">{item.problemTitle}</a>
                      {#if item.difficulty != null}
                        <div class="flex flex-row">
                          <span class={`rounded-full ${bgColorStyles.get(item.color ?? "black")} px-2 py-0.5 text-xs text-white`}>{item.difficulty}</span>
                        </div>
                      {/if}
                    </div>
                  {/each}
                </div>
              {/if}
            {/if}
          </div>
        {/each}
      {/if}
    </div>
  </div>
</div>
