<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import type { NumericRange } from "$lib/request";
  import type { SearchUserFacet } from "$lib/response";

  export let facet: SearchUserFacet | null;

  let countries: string[] = $page.url.searchParams.getAll("country");
  let ratingRange: string | null;
  let birthYearRange: string | null;
  let joinCountRange: string | null;

  function parseRange(label: string): NumericRange {
    const [begin, end] = label.split("~");
    const result: NumericRange = { begin: null, end: null };
    if (begin.trim() !== "") {
      result.begin = Number(begin.trim());
    }
    if (end.trim() !== "") {
      result.end = Number(end.trim());
    }
    return result;
  }

  function filter() {
    const p = new URLSearchParams($page.url.searchParams);

    if (countries.length > 0) {
      p.delete("country");
      for (const c of countries) {
        p.append("country", c);
      }
    } else {
      p.delete("country");
    }

    if (ratingRange != null) {
      const range = parseRange(ratingRange);
      if (range.begin != null) {
        p.set("ratingFrom", range.begin.toString());
      } else {
        p.delete("ratingFrom");
      }
      if (range.end != null) {
        p.set("ratingTo", range.end.toString());
      } else {
        p.delete("ratingTo");
      }
    } else {
      p.delete("ratingFrom");
      p.delete("ratingTo");
    }

    if (birthYearRange != null) {
      const range = parseRange(birthYearRange);
      if (range.begin != null) {
        p.set("birthYearFrom", range.begin.toString());
      } else {
        p.delete("birthYearFrom");
      }
      if (range.end != null) {
        p.set("birthYearTo", range.end.toString());
      } else {
        p.delete("birthYearTo");
      }
    } else {
      p.delete("birthYearFrom");
      p.delete("birthYearTo");
    }

    if (joinCountRange != null) {
      const range = parseRange(joinCountRange);
      if (range.begin != null) {
        p.set("joinCountFrom", range.begin.toString());
      } else {
        p.delete("joinCountFrom");
      }
      if (range.end != null) {
        p.set("joinCountTo", range.end.toString());
      } else {
        p.delete("joinCountTo");
      }
    } else {
      p.delete("joinCountFrom");
      p.delete("joinCountTo");
    }

    goto(`/user?${p.toString()}`, { replaceState: false, noScroll: true, invalidateAll: true });
  }
</script>

<div class="flex flex-col items-center bg-white px-2 py-2 text-sm">
  {#if facet != null}
    {#if facet.country != null}
      <div class="my-1 w-full rounded-md px-2 py-1 shadow-sm shadow-gray-500">
        <p class="text-md font-semibold">国</p>
        <div class="flex flex-wrap gap-1">
          {#each facet.country as c (c.label)}
            <div class="flex flex-row items-center justify-between gap-1 px-1">
              <input
                id={c.label}
                type="checkbox"
                class=""
                bind:group={countries}
                value={c.label}
                on:change={() => {
                  filter();
                }}
              />
              <label for={c.label} class="text-sm">
                {c.label}
              </label>
              <label for={c.label} class="flex-grow text-right text-sm">
                {c.count}
              </label>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <div class="my-1 w-full rounded-md px-2 py-2 shadow-sm shadow-gray-500">
      <div class="flex flex-row justify-between">
        <p class="text-md font-semibold">誕生年</p>
        <button
          class="rounded-full bg-gray-500 px-2 py-1 text-sm text-white"
          on:click={() => {
            birthYearRange = null;
            filter();
          }}>選択解除</button
        >
      </div>
      {#if facet.birthYear != null}
        {#each facet.birthYear as c (c.label)}
          <div class="flex-rows my-0.5 flex w-full items-center justify-between gap-2 px-3">
            <input
              id={c.label}
              type="radio"
              class=""
              bind:group={birthYearRange}
              value={c.label}
              on:change={() => {
                filter();
              }}
            />
            <label for={c.label} class="flex-grow text-sm">
              {c.label}
            </label>
            <label for={c.label} class="text-sm">
              {c.count}
            </label>
          </div>
        {/each}
      {/if}
    </div>

    <div class="my-1 w-full rounded-md px-2 py-2 shadow-sm shadow-gray-500">
      <div class="flex flex-row justify-between">
        <p class="text-md font-semibold">参加回数</p>
        <button
          class="rounded-full bg-gray-500 px-2 py-1 text-sm text-white"
          on:click={() => {
            joinCountRange = null;
            filter();
          }}>選択解除</button
        >
      </div>
      {#if facet.joinCount != null}
        {#each facet.joinCount as c (c.label)}
          <div class="flex-rows my-0.5 flex w-full items-center justify-between gap-2 px-3">
            <input
              id={c.label}
              type="radio"
              class=""
              bind:group={joinCountRange}
              value={c.label}
              on:change={() => {
                filter();
              }}
            />
            <label for={c.label} class="flex-grow text-sm">
              {c.label}
            </label>
            <label for={c.label} class="text-sm">
              {c.count}
            </label>
          </div>
        {/each}
      {/if}
    </div>
  {/if}
</div>
