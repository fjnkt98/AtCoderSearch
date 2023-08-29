<script lang="ts">
  import type { ProblemFacet, FilterRange } from "$lib/search";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";

  export let facet: ProblemFacet;

  const params = new URLSearchParams($page.url.searchParams);
  let categories: string[] = params.get("filter.category")?.split(",") ?? [];
  let difficultySelection: string = "";
  let difficulty = {
    from: params.get("filter.difficulty.from"),
    to: params.get("filter.difficulty.to"),
  };

  function filter() {
    const params = new URLSearchParams($page.url.searchParams);
    params.set("page", "1");

    if (categories.length > 0) {
      params.set("filter.category", categories.join(","));
    } else {
      params.delete("filter.category");
    }

    if (difficulty.from != null && difficulty.from !== "") {
      params.set("filter.difficulty.from", difficulty.from);
    } else {
      params.delete("filter.difficulty.from");
    }
    if (difficulty.to != null && difficulty.to !== "") {
      params.set("filter.difficulty.to", difficulty.to);
    } else {
      params.delete("filter.difficulty.to");
    }

    goto(`${$page.url.pathname}?${params.toString()}`, { replaceState: false });
  }
</script>

<div class="rounded-xl bg-white px-6 py-2 shadow-md shadow-gray-300">
  <button
    class="mb-2 mt-4 w-full rounded-xl bg-gray-600 px-2 py-1 text-white"
    on:click={() => {
      categories = [];
      difficultySelection = "";
      difficulty = { from: null, to: null };
      filter();
    }}>Reset</button
  >

  {#if facet.category != null}
    <p class="my-1 text-lg">Category</p>
    {#each facet.category as part}
      <div class="flex flex-row items-center justify-between">
        <input
          id={part.label}
          type="checkbox"
          bind:group={categories}
          on:change={() => {
            filter();
          }}
          value={part.label}
          class="mx-1 inline-block cursor-pointer"
        />
        <label for={part.label} class="mx-1 inline-block flex-grow cursor-pointer whitespace-nowrap text-left">
          {part.label}
        </label>
        <label for={part.label} class="mx-1 inline-block cursor-pointer">
          {part.count}
        </label>
      </div>
    {/each}
  {/if}

  {#if facet.difficulty != null}
    <p class="mt-4 text-lg">Difficulty</p>
    {#each facet.difficulty as part}
      <div class="flex flex-row items-center justify-between">
        <input
          id={part.label}
          type="radio"
          class="mx-1 inline-block cursor-pointer"
          bind:group={difficultySelection}
          on:change={() => {
            const splitted = part.label.split("~");
            if (splitted.length === 2) {
              const [from, to] = splitted;
              difficulty = {
                from: from.trim(),
                to: to.trim(),
              };
            } else if (splitted.length === 1) {
              if (part.label.startsWith("~")) {
                difficulty = {
                  from: null,
                  to: splitted[0].trim(),
                };
              } else if (part.label.endsWith("~")) {
                difficulty = {
                  from: splitted[0].trim(),
                  to: null,
                };
              }
            } else {
              difficulty = {
                from: null,
                to: null,
              };
            }
            filter();
          }}
          value={part.label}
        />
        <label for={part.label} class="mx-1 inline-block flex-grow cursor-pointer whitespace-nowrap text-left">
          {part.label}
        </label>
        <label for={part.label} class="mx-1 inline-block cursor-pointer">
          {part.count}
        </label>
      </div>
    {/each}
  {/if}
</div>
