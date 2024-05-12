<script lang="ts">
  import Icon from "svelte-icons-pack/Icon.svelte";
  import AiOutlineSearch from "svelte-icons-pack/ai/AiOutlineSearch";
  import HiOutlinePaperAirplane from "svelte-icons-pack/hi/HiOutlinePaperAirplane";
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher<{ search: string }>();

  export let q: string = "";
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
        dispatch("search", q);
      }
    }}
    bind:value={q}
  />
  <button
    type="button"
    class=""
    on:click|preventDefault={() => {
      dispatch("search", q);
    }}
  >
    <Icon src={HiOutlinePaperAirplane} className="rotate-90" />
  </button>
</div>
