<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import PageNavigation from "$lib/PageNavigation.svelte";
  import type { SubmissionResult } from "$lib/search";
  import Filter from "./Filter.svelte";
  import Submission from "./Submission.svelte";
  import type { Data } from "./data";

  export let data: Data;

  const labels = new Map<string, string>([
    ["-submitted_at", "提出日時降順"],
    ["submitted_at", "提出日時昇順"],
    ["-length", "コード長降順"],
    ["length", "コード長昇順"],
    ["-execution_time", "実行時間降順"],
    ["execution_time", "実行時間昇順"],
  ]);

  let selections: string[] = ["-submitted_at", "submitted_at", "-length", "length", "-execution_time", "execution_time"];
  let selected = $page.url.searchParams.get("sort") ?? "-submitted_at";

  let expand: boolean = false;
</script>

<div class="w-full flex-1 overflow-y-auto py-8 sm:px-8">
  <PageNavigation end={data.result.stats.pages} current={data.result.stats.index} />

  <div class="my-2 flex min-w-min flex-row items-start justify-center">
    <div class={`mx-4 ${expand ? "block basis-1/5" : "hidden"} lg:block`}>
      <select
        class="my-2 block w-full rounded-lg bg-white p-2.5 text-sm shadow-sm shadow-gray-300"
        bind:value={selected}
        on:change={() => {
          const params = new URLSearchParams($page.url.searchParams);
          params.set("sort", selected);
          params.set("page", "1");
          goto(`${$page.url.pathname}?${params.toString()}`);
        }}
      >
        {#each selections as s}
          <option value={s}>
            {labels.get(s)}
          </option>
        {/each}
      </select>

      <Filter categories={data.categories} languages={data.languages} contests={data.contests} problems={data.problems} />
    </div>

    <div class={`mx-4 flex ${expand ? "basis-4/5" : "sm:basis-4/5"} flex-col items-center justify-center`}>
      <div class="flex-rows flex w-full max-w-5xl items-center sm:justify-between md:justify-between lg:justify-end">
        <button
          class={`rounded-xl ${expand ? "bg-green-600 text-slate-50" : "bg-white text-green-600"} px-4 py-1 text-lg font-medium shadow-sm shadow-gray-300 lg:hidden`}
          on:click={() => {
            expand = !expand;
          }}>Filter</button
        >
        <p class="my-2 w-full max-w-5xl text-right text-slate-500">{data.result.stats.count}件/{data.result.stats.total}件 約{data.result.stats.time}ms</p>
      </div>
      {#each data.result.items as item}
        <Submission submission={item} />
      {/each}
    </div>
  </div>
</div>
