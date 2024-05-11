<script lang="ts">
  import { page } from "$app/stores";
  import Icon from "svelte-icons-pack/Icon.svelte";
  import AiOutlineSearch from "svelte-icons-pack/ai/AiOutlineSearch";
  import HiOutlinePaperAirplane from "svelte-icons-pack/hi/HiOutlinePaperAirplane";
  import { goto } from "$app/navigation";

  export let href: string;
  export let s: string;

  let q: string = $page.url.searchParams.get("q") ?? "";

  function search() {
    const p = new URLSearchParams();
    p.set("q", q);
    p.set("s", s);
    goto(`${href}?${p.toString()}`, { replaceState: true, invalidateAll: true });
  }
</script>

<div class="flex flex-row items-center justify-center rounded-full border border-gray-500 bg-white p-1 text-gray-900">
  <div class="px-2">
    <Icon src={AiOutlineSearch} />
  </div>
  <input
    type="search"
    class="w-full appearance-none bg-transparent text-gray-900 focus:border-none focus:outline-none"
    placeholder="Search"
    on:keydown={(e) => {
      if (e.key === "Enter") {
        e.preventDefault();
        search();
      }
    }}
    bind:value={q}
  />
  <button
    type="button"
    class=""
    on:click|preventDefault={() => {
      search();
    }}
  >
    <Icon src={HiOutlinePaperAirplane} className="rotate-90" />
  </button>
</div>
