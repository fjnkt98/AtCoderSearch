<script lang="ts">
  import type { ProblemFacet } from "$lib/search";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";

  export let facet: ProblemFacet;

  const params = new URLSearchParams($page.url.searchParams);
  let categories: string[] = params.get("filter.category")?.split(",") ?? [];
  let colors: string[] = params.get("filter.color")?.split(",") ?? [];

  function filter() {
    const params = new URLSearchParams($page.url.searchParams);
    params.set("page", "1");
    params.set("filter.category", categories.join(","));
    params.set("filter.color", colors.join(","));

    goto(`${$page.url.pathname}?${params.toString()}`, { replaceState: false });
  }
</script>

<div class="rounded-xl bg-white px-6 py-2 shadow-md shadow-gray-300">
  <p class="my-1 text-lg">Category</p>
  {#each facet.category as part}
    <div class="flex flex-row items-center justify-between">
      <input
        id={part.label}
        type="checkbox"
        bind:group={categories}
        value={part.label}
        on:change={() => {
          filter();
        }}
        class="mx-1 inline-block cursor-pointer"
      />
      <label for={part.label} class="mx-1 inline-block flex-grow cursor-pointer text-left">
        {part.label}
      </label>
      <label for={part.label} class="mx-1 inline-block cursor-pointer">
        {part.count}
      </label>
    </div>
  {/each}

  <p class="text-lg">Color</p>
  {#each facet.color as part}
    <div class="flex flex-row items-center justify-between">
      <input
        id={part.label}
        type="checkbox"
        bind:group={colors}
        value={part.label}
        on:change={() => {
          filter();
        }}
        class="mx-1 inline-block cursor-pointer"
      />
      <label for={part.label} class="mx-1 inline-block flex-grow cursor-pointer text-left">
        {part.label}
      </label>
      <label for={part.label} class="mx-1 inline-block cursor-pointer">
        {part.count}
      </label>
    </div>
  {/each}

  <button
    class="mb-2 mt-4 w-full rounded-xl bg-gray-600 px-2 py-1 text-white"
    on:click={() => {
      categories = [];
      colors = [];
      filter();
    }}>Reset</button
  >
</div>
