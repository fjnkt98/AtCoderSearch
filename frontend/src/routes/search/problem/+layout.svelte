<script lang="ts">
  import SearchBar from "$lib/SearchBar.svelte";
  import { page } from "$app/stores";
  import { defaultProblemSearchParams } from "$lib/defaults";

  const keyword: string = $page.url.searchParams.get("keyword") ?? "";
</script>

<nav class="flex w-full items-center justify-between border px-6 py-1">
  <a href="/" class="font-roboto mx-4 my-2 text-3xl">AtCoder Search</a>
  <SearchBar path={"/search/problem"} defaultParams={defaultProblemSearchParams} {keyword} />
</nav>

<slot />
