<script lang="ts">
  import { goto } from "$app/navigation";
  // @ts-ignore
  import Icon from "svelte-icons-pack/Icon.svelte";
  import AiOutlineSearch from "svelte-icons-pack/ai/AiOutlineSearch";
  import HiOutlinePaperAirplane from "svelte-icons-pack/hi/HiOutlinePaperAirplane";

  export let path: string;
  export let defaultParams: URLSearchParams;

  export let keyword: string = "";

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

<div class="flex flex-row items-center justify-center rounded-full border border-gray-400 bg-white p-1 text-gray-900 sm:p-2">
  <div class="px-2">
    <Icon src={AiOutlineSearch} />
  </div>
  <input type="search" class="w-full appearance-none bg-transparent text-gray-900 focus:border-none focus:outline-none" placeholder="Search" on:keydown={handleKeyDown} bind:value={keyword} />
  <button
    type="button"
    class=""
    on:click={() => {
      search();
    }}
  >
    <Icon src={HiOutlinePaperAirplane} className="rotate-90" />
  </button>
</div>
