<script lang="ts">
  import SearchBar from "$lib/SearchBar.svelte";
  import { page } from "$app/stores";

  const keyword: string = $page.url.searchParams.get("keyword") ?? "";
  const defaultParams = new URLSearchParams([
    ["limit", "20"],
    ["page", "1"],
    ["facet.term", "category"],
    ["facet.difficulty.from", "0"],
    ["facet.difficulty.to", "4000"],
    ["facet.difficulty.gap", "400"],
  ]);
</script>

<nav class="flex w-full items-center justify-between border px-6 py-1">
  <a href="/" class="mx-4 my-2 font-roboto text-3xl">AtCoder Search</a>
  <SearchBar path={"/search/problem"} {defaultParams} {keyword} />
</nav>

<slot />
