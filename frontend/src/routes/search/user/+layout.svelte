<script lang="ts">
  import SearchBar from "$lib/SearchBar.svelte";
  import { page } from "$app/stores";
  import { defaultUserSearchParams } from "$lib/defaults";

  const keyword: string = $page.url.searchParams.get("keyword") ?? "";
</script>

<nav class="flex w-full items-center justify-between border px-2 py-1 sm:px-6">
  <a href="/" class="mx-4 my-2 font-roboto text-xl sm:text-3xl">AtCoder Search</a>
  <div class="basis-1/2">
    <SearchBar path={"/search/user"} defaultParams={defaultUserSearchParams} {keyword} />
  </div>
</nav>

<slot />
