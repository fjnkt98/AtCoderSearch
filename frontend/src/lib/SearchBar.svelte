<script lang="ts">
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import Icon from "svelte-icons-pack/Icon.svelte";
  import AiOutlineSearch from "svelte-icons-pack/ai/AiOutlineSearch";

  export let path: string;
  export let defaultParams: URLSearchParams;

  let keyword: string = "";

  function search() {
    const params = new URLSearchParams(defaultParams);
    params.set("keyword", keyword);
    goto(`${path}?${params.toString()}`, { replaceState: false });
  }

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      e.preventDefault();
      search();
    }
  }
</script>

<div class="white flex flex-row items-center justify-center rounded-full border border-gray-400 p-2 text-gray-900">
  <div class="px-2">
    <Icon src={AiOutlineSearch} />
  </div>
  <input type="text" class="flex-1 appearance-none bg-transparent text-gray-900 focus:border-none focus:outline-none" placeholder="Search" on:keydown={handleKeyDown} bind:value={keyword} />
  <button
    type="button"
    class="bg-transparent"
    on:click={() => {
      search();
    }}
  />
</div>
