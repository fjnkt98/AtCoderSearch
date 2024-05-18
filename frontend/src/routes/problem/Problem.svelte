<script lang="ts">
  import { bgColorStyles, textColorStyles } from "$lib/colors";
  import type { Problem } from "$lib/response";
  import { createQuery } from "@tanstack/svelte-query";
  import dayjs from "dayjs";
  import timezone from "dayjs/plugin/timezone";
  import utc from "dayjs/plugin/utc";
  import RecommendItem from "./RecommendItem.svelte";
  import { fetchRecommendProblemResult } from "./query";

  dayjs.extend(timezone);
  dayjs.extend(utc);

  export let item: Problem;

  let showDetail: boolean = false;

  $: recommendQuery = createQuery({
    queryKey: [`recommendTo${item.problemId}`, showDetail],
    queryFn: async () => {
      const res = await fetchRecommendProblemResult({ limit: 5, page: 1, problemId: item.problemId }, fetch);
      return res;
    },
    enabled: showDetail,
  });

  function convertDateTime(date: string): string {
    return dayjs(date).tz("Asia/Tokyo").format("YYYY/MM/DD HH:mm:ss");
  }
</script>

<div class="mx-2 my-1 flex flex-col rounded-lg bg-white px-6 py-3 shadow-sm shadow-gray-400">
  <div class="flex flex-row items-center justify-between">
    <a href={item.contestUrl} target="_blank" class="block font-medium text-blue-500">
      {item.contestId.toUpperCase()}
    </a>
    <a href={item.problemUrl} target="_blank" class={`block text-pretty font-medium ${textColorStyles.get(item.color ?? "black")}`}>
      {item.problemTitle}
    </a>
    <label class={`w-12 cursor-pointer text-sm text-black ` + (showDetail ? "" : "underline")}>
      {showDetail ? "閉じる" : "詳細"}
      <input
        type="button"
        class="h-0 w-0"
        on:click={() => {
          showDetail = !showDetail;
        }}
      />
    </label>
  </div>
  {#if showDetail}
    <div class="flex flex-row py-1">
      <div class="my-1 flex flex-grow flex-col items-start justify-between">
        <span class="block text-sm text-slate-500">{item.contestTitle}</span>
        <span class="block text-sm text-slate-500">{convertDateTime(item.startAt)}</span>
      </div>
      <div class="flex flex-col items-end justify-end gap-1">
        <span class={`block rounded-full px-2 py-0.5 text-xs font-medium text-white  ${bgColorStyles.get(item.category)}`}>{item.category}</span>
        {#if item.difficulty != null}
          <span class={`block rounded-full px-2 py-0.5 text-xs font-medium text-white ${bgColorStyles.get(item.color ?? "black")}`}>{item.difficulty}</span>
        {/if}
      </div>
    </div>

    {#if $recommendQuery.isSuccess}
      <p class="mt-2 text-gray-700">似ているかも</p>
      <div class="my-1 flex flex-row gap-4 overflow-x-scroll py-0.5">
        {#each $recommendQuery.data.items as item}
          <RecommendItem {item} />
        {/each}
      </div>
    {/if}
  {/if}
</div>
