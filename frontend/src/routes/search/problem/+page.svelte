<script lang="ts">
  import type { ProblemSearchResult } from "$lib/search";

  export let data: ProblemSearchResult;
</script>

<div class="flex-1 overflow-auto p-16">
  {#await data.items}
    <p>search...</p>
  {:then items}
    {#each items as item}
      <p>{item.problem_title}</p>
    {/each}
  {:catch}
    <p>Search failed!</p>
  {/await}
</div>
