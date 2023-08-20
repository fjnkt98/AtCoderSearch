<script lang="ts">
  import SearchBar from "$lib/SearchBar.svelte";

  const defaultParams = new URLSearchParams({
    limit: "20",
    page: "1",
    facet: "category,color",
  });
</script>

<nav class="flex w-full items-center justify-between border px-6 py-1">
  <a href="/" class="mx-4 my-2 font-roboto text-3xl">AtCoder Search</a>
  <SearchBar path={"/search/problem"} {defaultParams} />
</nav>

<slot />
