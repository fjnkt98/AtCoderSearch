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

  let expand: boolean = false;
</script>

<div class="w-full flex-1 overflow-y-auto py-8 sm:px-8">
  <PageNavigation end={data.stats.pages} current={data.stats.index} />

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
      {#if data.stats.facet != null}
        <Facet facet={data.stats.facet} />
      {/if}
    </div>

    <div class={`mx-4 flex ${expand ? "basis-4/5" : "sm:basis-4/5"} flex-col items-center justify-center`}>
      <div class="flex-rows flex w-full max-w-5xl items-center sm:justify-between md:justify-between lg:justify-end">
        <button
          class={`rounded-xl ${expand ? "bg-green-600 text-slate-50" : "bg-white text-green-600"} px-4 py-1 text-lg font-medium shadow-sm shadow-gray-300 lg:hidden`}
          on:click={() => {
            expand = !expand;
          }}>Filter</button
        >
        <p class="my-2 w-full max-w-5xl text-right text-slate-500">{data.stats.count}件/{data.stats.total}件 約{data.stats.time}ms</p>
      </div>
      {#each data.items as item}
        <Problem problem={item} />
      {/each}
    </div>
  </div>
</div>
