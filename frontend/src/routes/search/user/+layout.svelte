<script lang="ts">
  import SearchBar from "$lib/SearchBar.svelte";
  import { page } from "$app/stores";

  const keyword: string = $page.url.searchParams.get("keyword") ?? "";
  const defaultParams = new URLSearchParams({
    limit: "50",
    page: "1",
    facet: "color,birth_year,country",
  });
</script>

<nav class="flex w-full items-center justify-between border px-6 py-1">
  <a href="/" class="mx-4 my-2 font-roboto text-3xl">AtCoder Search</a>
  <SearchBar path={"/search/user"} {defaultParams} {keyword} />
</nav>

<slot />
