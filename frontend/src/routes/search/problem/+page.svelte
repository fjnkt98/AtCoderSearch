<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import PageNavigation from "$lib/PageNavigation.svelte";
  import type { ProblemSearchResult } from "$lib/search";
  import Facet from "./Facet.svelte";
  import Problem from "./Problem.svelte";

  export let data: ProblemSearchResult;

  const labels = new Map<string, string>([
    ["-score", "検索スコア順"],
    ["start_at", "開催日時早い順"],
    ["-start_at", "開催日時遅い順"],
    ["difficulty", "難易度低い順"],
    ["-difficulty", "難易度高い順"],
  ]);

  let selections: string[] = ["-start_at", "start_at", "-score", "difficulty", "-difficulty"];
  let selected = $page.url.searchParams.get("sort") ?? "-start_at";
</script>

<div class="flex-1 overflow-auto px-12 py-8">
  <PageNavigation end={data.stats.pages} current={data.stats.index} />

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
      {#if data.stats.facet != null}
        <Facet facet={data.stats.facet} />
      {/if}
    </div>

    <div class="mx-2 flex flex-1 flex-col items-center justify-center">
      <p class="my-2 w-2/3 min-w-[600px] text-left text-slate-500">{data.stats.count}件/{data.stats.total}件 約{data.stats.time}ms</p>
      {#each data.items as item}
        <Problem problem={item} />
      {/each}
    </div>
  </div>
</div>
