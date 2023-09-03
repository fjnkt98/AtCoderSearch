<script lang="ts">
  import { goto } from "$app/navigation";
  import Icon from "svelte-icons-pack/Icon.svelte";
  import { defaultRecommendParams } from "$lib/defaults";
  import AiOutlineSearch from "svelte-icons-pack/ai/AiOutlineSearch";
  import HiOutlinePaperAirplane from "svelte-icons-pack/hi/HiOutlinePaperAirplane";

  export let data;
  let userId: string = data.userId ?? "";

  function search() {
    const params = new URLSearchParams(defaultRecommendParams);
    params.set("user_id", userId);
    goto(`/recommend/problem?${params.toString()}`, { replaceState: false });
  }

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      e.preventDefault();
      search();
    }
  }
</script>

<nav class="flex w-full items-center justify-between border px-6 py-1">
  <a href="/" class="mx-4 my-2 font-roboto text-3xl">AtCoder Search</a>
  <div class="white flex flex-row items-center justify-center rounded-full border border-gray-400 bg-white p-2 text-gray-900">
    <div class="px-2">
      <Icon src={AiOutlineSearch} />
    </div>
    <input type="text" class="flex-1 appearance-none bg-transparent text-gray-900 focus:border-none focus:outline-none" placeholder="User ID" on:keydown={handleKeyDown} bind:value={userId} />
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
</nav>

<slot />
