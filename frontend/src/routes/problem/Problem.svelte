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

<div class="mx-2 my-1 flex min-w-0 flex-col rounded-lg bg-white px-2 py-1 shadow-sm shadow-gray-400 lg:min-w-96">
  <div class="my-1.5 flex flex-row items-center justify-between text-sm md:mx-4 md:text-base">
    <a href={item.contestUrl} target="_blank" class="block text-wrap font-medium text-blue-500">
      {item.contestId.toUpperCase()}
    </a>
    <a href={item.problemUrl} target="_blank" class={`mx-1.5 block w-1/2 text-center font-medium ${textColorStyles.get(item.color ?? "black")}`}>
      {item.problemTitle}
    </a>
    <div class="flex flex-row items-center gap-1 text-xs">
      {#if item.difficulty != null}
        <span class={`block rounded-xl px-1.5 py-0.5 text-xs font-medium text-white ${bgColorStyles.get(item.color ?? "black")}`}>{item.difficulty}</span>
      {/if}
      <span class={`block text-wrap rounded-xl  px-1.5 py-0.5 text-center text-xs font-medium text-white  ${bgColorStyles.get(item.category)}`}>{item.category}</span>
      <label class={`w-12 cursor-pointer text-center text-sm text-black ` + (showDetail ? "" : "underline")}>
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
  </div>
  {#if showDetail}
    <div class="mx-3 flex flex-row py-1">
      <div class="my-1 flex flex-grow flex-col items-start justify-between">
        <span class="block text-xs text-slate-500">{item.contestTitle}</span>
        <span class="block text-xs text-slate-500">{convertDateTime(item.startAt)}</span>
      </div>
    </div>

    {#if $recommendQuery.isSuccess}
      <p class="mx-2 mt-2 text-sm text-gray-700">似ているかも</p>
      <div class="my-1 flex flex-row gap-4 overflow-x-scroll py-0.5">
        {#each $recommendQuery.data.items as item}
          <RecommendItem {item} />
        {/each}
      </div>
    {/if}
  {/if}
</div>
