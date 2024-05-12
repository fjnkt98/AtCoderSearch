<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import Header from "$lib/Header.svelte";
  import FilterMenu from "./FilterMenu.svelte";
  import SearchBar from "$lib/SearchBar.svelte";
  import { createQuery } from "@tanstack/svelte-query";
  import Tab from "$lib/Tab.svelte";
  import { bgColorStyles, textColorStyles } from "$lib/colors";
  import type { SearchProblemResult } from "$lib/response";
  import { fetchSearchProblemResult } from "./query";
  import dayjs from "dayjs";
  import timezone from "dayjs/plugin/timezone";
  import utc from "dayjs/plugin/utc";
  import { selections } from "./sort";

  dayjs.extend(timezone);
  dayjs.extend(utc);

  $: width = 0;
  const widthThreshold = 768;

  let opened = false;
  let detail: string | null = null;
  let s: string = $page.url.searchParams.get("s") ?? "2";

  export let data: SearchProblemResult;
  const searchQuery = createQuery({
    queryKey: ["search"],
    queryFn: async () => {
      const res = await fetchSearchProblemResult($page.url.searchParams, fetch);
      return res;
    },
    initialData: data,
  });

  function convertDateTime(date: string): string {
    return dayjs(date).tz("Asia/Tokyo").format("YYYY/MM/DD HH:mm:ss");
  }
</script>

<svelte:window bind:outerWidth={width} />

<div class="flex h-dvh flex-col bg-gray-100">
  <Header />
  <Tab selected={"problem"} />

  <div class={`flex ${width > widthThreshold ? "flex-row justify-between" : "flex-col"}`}>
    <div class="mx-2 my-1">
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

    <div class="flex-grow">
      <div class="my-1 flex flex-row items-center justify-between gap-2 px-2 text-sm">
        <select
          class="block rounded-md border border-gray-400 bg-white px-2 py-1"
          bind:value={s}
          on:change={() => {
            const p = $page.url.searchParams;
            p.set("s", s.toString());
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
          <SearchBar href={"/problem"} s={"1"} />
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
              <label class={`text-sm text-black ` + (detail == item.problemId ? "" : "underline")}>
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
              <div class="py-1">
                <div class="my-0.5 flex flex-row items-end justify-between">
                  <span class="block text-sm text-slate-500">{item.contestTitle}</span>
                  <span class={`block rounded-full px-2 py-0.5 text-xs font-medium text-white  ${bgColorStyles.get(item.category)}`}>{item.category}</span>
                </div>
                <div class="flex flex-row items-center justify-between">
                  <span class="block text-sm text-slate-500">{convertDateTime(item.startAt)}</span>
                  {#if item.difficulty != null}
                    <span class={`block rounded-full px-2 py-0.5 text-xs font-medium text-white ${bgColorStyles.get(item.color ?? "black")}`}>{item.difficulty}</span>
                  {/if}
                </div>
              </div>
            {/if}
          </div>
        {/each}
      {/if}
    </div>
  </div>
</div>
