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
</script>

<div class="flex-1 overflow-auto px-12 py-8">
  <PageNavigation end={data.result.stats.pages} current={data.result.stats.index} />

  <div class="container flex flex-row justify-between">
    <div>
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

    <div class="mx-2 flex flex-1 flex-col items-center justify-center">
      <p class="my-2 w-2/3 min-w-[600px] text-left text-slate-500">{data.result.stats.count}件/{data.result.stats.total}件 約{data.result.stats.time}ms</p>
      {#each data.result.items as item}
        <Submission submission={item} />
      {/each}
    </div>
  </div>
</div>
