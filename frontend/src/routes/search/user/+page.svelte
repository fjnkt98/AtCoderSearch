<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import PageNavigation from "$lib/PageNavigation.svelte";
  import type { UserSearchResult } from "$lib/search";
  import Facet from "./Facet.svelte";
  import User from "./User.svelte";

  export let data: UserSearchResult;

  const labels = new Map<string, string>([
    ["-rating", "レート高い順"],
    ["rating", "レート低い順"],
    ["birth_year", "誕生日早い順"],
    ["-birth_year", "誕生日遅い順"],
    ["-score", "検索スコア順"],
  ]);

  let selections: string[] = ["-rating", "rating", "birth_year", "-birth_year", "-score"];
  let selected = $page.url.searchParams.get("sort") ?? "-rating";
</script>

<div class="flex-1 overflow-auto px-12 py-8">
  <PageNavigation end={data.stats.pages} current={data.stats.index} />

  <div class="container flex flex-row items-start justify-between">
    <div class="basis-1/6">
      <select
        class="my-2 block w-full rounded-lg bg-white p-2.5 text-sm shadow-sm shadow-gray-300"
        bind:value={selected}
        on:change={() => {
          const params = new URLSearchParams($page.url.searchParams);
          params.set("sort", selected);
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
        <User user={item} />
      {/each}
    </div>
  </div>
</div>
