<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import type { NumericRange } from "$lib/request";
  import type { SearchProblemFacet } from "$lib/response";

  export let facet: SearchProblemFacet | null;

  let categories: string[] = [];
  let difficultyRange: string | null;

  let difficulty: number | null;
  let userId: string | null = $page.url.searchParams.get("userId");
  let excludeSolved: boolean = false;
  let excludeExperimental: boolean = false;
  let prioritizeRecent: boolean = false;

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

    if (categories.length > 0) {
      p.delete("category");
      for (const c of categories) {
        p.append("category", c);
      }
    } else {
      p.delete("category");
    }

    if (difficulty != null) {
      p.set("difficulty", difficulty.toString());
    } else {
      p.delete("difficulty");
    }

    if (difficultyRange != null) {
      const range = parseRange(difficultyRange);
      if (range.begin != null) {
        p.set("difficultyFrom", range.begin.toString());
      } else {
        p.delete("difficultyFrom");
      }
      if (range.end != null) {
        p.set("difficultyTo", range.end.toString());
      } else {
        p.delete("difficultyTo");
      }
    } else {
      p.delete("difficultyFrom");
      p.delete("difficultyTo");
    }

    if (userId != null && userId !== "" && excludeSolved) {
      p.set("userId", userId);
      p.set("excludeSolved", "true");
    } else {
      p.delete("userId");
      p.delete("excludeSolved");
    }

    if (excludeExperimental) {
      p.set("experimental", "false");
    } else {
      p.delete("experimental");
    }

    if (prioritizeRecent) {
      p.set("prioritizeRecent", "true");
    } else {
      p.delete("prioritizeRecent");
    }

    goto(`/problem?${p.toString()}`, { replaceState: false, noScroll: true, invalidateAll: true });
  }
</script>

<div class="flex flex-col items-center bg-white px-2 py-2 text-sm">
  {#if facet != null}
    {#if facet.category != null}
      <div class="my-1 w-full rounded-md px-2 py-1 shadow-sm shadow-gray-500">
        <p class="text-md font-semibold">カテゴリ</p>
        {#each facet.category as c (c.label)}
          <div class="flex-rows my-0.5 flex w-full items-center justify-between gap-2 px-3">
            <input
              id={c.label}
              type="checkbox"
              class=""
              bind:group={categories}
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
      </div>
    {/if}

    <div class="my-1 w-full rounded-md px-2 py-2 shadow-sm shadow-gray-500">
      <div class="flex flex-row justify-between">
        <p class="text-md font-semibold">難易度</p>
        <button
          class="rounded-full bg-gray-500 px-2 py-1 text-sm text-white"
          on:click={() => {
            difficultyRange = null;
            filter();
          }}>選択解除</button
        >
      </div>
      <div class="">
        <input
          type="number"
          class="mx-2 mb-1 w-1/2 rounded-md border border-gray-600 bg-transparent px-2 text-gray-700 focus:border-blue-700"
          placeholder="近い難易度の問題を探す"
          bind:value={difficulty}
          on:keydown={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();

              difficultyRange = null;
              filter();
            }
          }}
        />
      </div>
      {#if facet.difficulty != null}
        {#each facet.difficulty as c (c.label)}
          <div class="flex-rows my-0.5 flex w-full items-center justify-between gap-2 px-3">
            <input
              id={c.label}
              type="radio"
              class=""
              bind:group={difficultyRange}
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

  <div class="flex w-full flex-col rounded-md px-2 py-2 shadow-sm shadow-gray-500">
    <p class="text-md font-semibold">その他</p>
    <input
      type="search"
      class="w-1/2 rounded-md border border-gray-900 bg-transparent px-2 text-gray-700 focus:border-blue-700"
      placeholder="ユーザID"
      bind:value={userId}
      on:keydown={(e) => {
        if (e.key === "Enter") {
          e.preventDefault();
          filter();
        }
      }}
    />

    <label class="">
      <input
        type="checkbox"
        class=""
        bind:checked={excludeSolved}
        on:change={() => {
          filter();
        }}
      />
      解いたことのある問題を除外
    </label>

    <label class="">
      <input
        type="checkbox"
        class=""
        bind:checked={excludeExperimental}
        on:change={() => {
          filter();
        }}
      />
      試験管問題を除外
    </label>

    <label class="">
      <input
        type="checkbox"
        class=""
        bind:checked={prioritizeRecent}
        on:change={() => {
          filter();
        }}
      />
      最近の問題を優先して表示 (検索スコア順でのみ)
    </label>
  </div>
</div>
