<script lang="ts">
  import type { UserFacet } from "$lib/search";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";

  export let facet: UserFacet;

  const params = new URLSearchParams($page.url.searchParams);
  let colors: string[] = params.get("filter.color")?.split(",") ?? [];
  let years: string[] = params.get("filter.birth_year")?.split(",") ?? [];
  let countries: string[] = params.get("filter.country")?.split(",") ?? [];

  function filter() {
    const params = new URLSearchParams($page.url.searchParams);
    params.set("page", "1");
    params.set("filter.color", colors.join(","));
    params.set("filter.birth_year", years.join(","));
    params.set("filter.country", countries.join(","));

    goto(`${$page.url.pathname}?${params.toString()}`, { replaceState: false });
  }
</script>

<div class="rounded-xl bg-white px-6 py-4 shadow-md shadow-gray-300">
  {#if facet.color != null}
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
  {/if}

  {#if facet.birth_year != null}
    <p class="text-lg">Birth Year</p>
    {#each facet.birth_year as part}
      <div class="flex flex-row items-center justify-between">
        <input
          id={part.label}
          type="checkbox"
          bind:group={years}
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
  {/if}

  <button
    class="mb-2 mt-4 w-full rounded-xl bg-gray-600 px-2 py-1 text-white"
    on:click={() => {
      colors = [];
      years = [];
      countries = [];
      filter();
    }}>Reset</button
  >
</div>
