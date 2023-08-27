<script lang="ts">
  import type { UserFacet, FilterRange } from "$lib/search";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";

  export let facet: UserFacet;

  const params = new URLSearchParams($page.url.searchParams);
  let countries: string[] = params.get("filter.country")?.split(",") ?? [];

  let ratingSelection: string = "";
  let rating: FilterRange = {
    from: params.get("filter.rating.from"),
    to: params.get("filter.rating.to"),
  };

  let birthYearSelection: string = "";
  let birthYear: FilterRange = {
    from: params.get("filter.birth_year.from"),
    to: params.get("filter.birth_year.to"),
  };

  let joinCountSelection: string = "";
  let joinCount: FilterRange = {
    from: params.get("filter.join_count.from"),
    to: params.get("filter.join_count.to"),
  };

  function filter() {
    const params = new URLSearchParams($page.url.searchParams);
    params.set("page", "1");

    if (countries.length > 0) {
      params.set("filter.country", countries.join(","));
    } else {
      params.delete("filter.country");
    }

    if (rating.from != null && rating.from !== "") {
      params.set("filter.rating.from", rating.from);
    } else {
      params.delete("filter.rating.from");
    }
    if (rating.to != null && rating.to !== "") {
      params.set("filter.rating.to", rating.to);
    } else {
      params.delete("filter.rating.to");
    }

    if (birthYear.from != null && birthYear.from !== "") {
      params.set("filter.birth_year.from", birthYear.from);
    } else {
      params.delete("filter.birth_year.from");
    }
    if (birthYear.to != null && birthYear.to !== "") {
      params.set("filter.birth_year.to", birthYear.to);
    } else {
      params.delete("filter.birth_year.to");
    }

    if (joinCount.from != null && joinCount.from !== "") {
      params.set("filter.join_count.from", joinCount.from);
    } else {
      params.delete("filter.join_count.from");
    }
    if (joinCount.to != null && joinCount.to !== "") {
      params.set("filter.join_count.to", joinCount.to);
    } else {
      params.delete("filter.join_count.to");
    }

    goto(`${$page.url.pathname}?${params.toString()}`, { replaceState: false });
  }
</script>

<div class="rounded-xl bg-white px-6 py-4 shadow-md shadow-gray-300">
  <button
    class="mb-2 mt-4 w-full rounded-xl bg-gray-600 px-2 py-1 text-white"
    on:click={() => {
      countries = [];
      ratingSelection = "";
      rating = { from: null, to: null };
      birthYearSelection = "";
      birthYear = { from: null, to: null };
      joinCountSelection = "";
      joinCount = { from: null, to: null };

      filter();
    }}>Reset</button
  >

  {#if facet.rating != null}
    <p class="text-lg">Rating</p>
    {#each facet.rating as part}
      <div class="flex flex-row items-center justify-between">
        <input
          id={part.label}
          type="radio"
          class="mx-1 inline-block cursor-pointer"
          bind:group={ratingSelection}
          value={part.label}
          on:change={() => {
            const splitted = part.label.split("~");
            if (splitted.length === 2) {
              const [from, to] = splitted;
              rating = {
                from: from.trim(),
                to: to.trim(),
              };
            } else if (splitted.length === 1) {
              if (part.label.startsWith("~")) {
                rating = {
                  from: null,
                  to: splitted[0].trim(),
                };
              } else if (part.label.endsWith("~")) {
                rating = {
                  from: splitted[0].trim(),
                  to: null,
                };
              }
            } else {
              rating = {
                from: null,
                to: null,
              };
            }
            filter();
          }}
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

  {#if facet.birth_year != null}
    <p class="text-lg">Birth Year</p>
    {#each facet.birth_year as part}
      <div class="flex flex-row items-center justify-between">
        <input
          id={part.label}
          type="radio"
          class="mx-1 inline-block cursor-pointer"
          bind:group={birthYearSelection}
          value={part.label}
          on:change={() => {
            const splitted = part.label.split("~");
            if (splitted.length === 2) {
              const [from, to] = splitted;
              birthYear = {
                from: from.trim(),
                to: to.trim(),
              };
            } else if (splitted.length === 1) {
              if (part.label.startsWith("~")) {
                birthYear = {
                  from: null,
                  to: splitted[0].trim(),
                };
              } else if (part.label.endsWith("~")) {
                birthYear = {
                  from: splitted[0].trim(),
                  to: null,
                };
              }
            } else {
              birthYear = {
                from: null,
                to: null,
              };
            }
            filter();
          }}
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

  {#if facet.join_count != null}
    <p class="text-lg">Join Count</p>
    {#each facet.join_count as part}
      <div class="flex flex-row items-center justify-between">
        <input
          id={part.label}
          type="radio"
          class="mx-1 inline-block cursor-pointer"
          bind:group={joinCountSelection}
          value={part.label}
          on:change={() => {
            const splitted = part.label.split("~");
            if (splitted.length === 2) {
              const [from, to] = splitted;
              joinCount = {
                from: from.trim(),
                to: to.trim(),
              };
            } else if (splitted.length === 1) {
              if (part.label.startsWith("~")) {
                joinCount = {
                  from: null,
                  to: splitted[0].trim(),
                };
              } else if (part.label.endsWith("~")) {
                joinCount = {
                  from: splitted[0].trim(),
                  to: null,
                };
              }
            } else {
              joinCount = {
                from: null,
                to: null,
              };
            }
            filter();
          }}
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

  {#if facet.country != null}
    <p class="my-1 text-lg">Country</p>
    {#each facet.country as part}
      <div class="flex flex-row items-center justify-between">
        <input
          id={part.label}
          type="checkbox"
          bind:group={countries}
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
</div>
