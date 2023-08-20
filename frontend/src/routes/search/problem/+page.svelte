<script lang="ts">
  import { page } from "$app/stores";
  import PageNavigation from "$lib/PageNavigation.svelte";
  import type { ProblemSearchResult } from "$lib/search";
  import Problem from "./Problem.svelte";

  export let data: ProblemSearchResult;

  const params = new URLSearchParams($page.url.searchParams);
  let current: number = data.stats.index;
  let pages: number = data.stats.pages;
</script>

<div class="flex-1 overflow-auto p-16">
  <PageNavigation end={pages} {current} path={"/search/problem"} {params} />

  <div class="container flex flex-row justify-between">
    <div class="border-2">Facets</div>
    <div class="flex flex-1 flex-col items-center justify-center">
      {#each data.items as item}
        <Problem problem={item} />
      {/each}
    </div>
  </div>
</div>
