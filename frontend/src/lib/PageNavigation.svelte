<script lang="ts">
  export let end: number;
  export let current: number;
  export let params: URLSearchParams;
  export let path: string;

  type Navigation = {
    label: string;
    path: string;
  };

  const width = 5;

  let navigation: Array<Navigation> = [];
  let count = 0;
  let containsBegin = false;
  let containsEnd = false;
  for (let i = Math.max(current - width, 1); i <= Math.min(current + width, end); i++) {
    params.set("page", i.toString());
    if (i === 1) {
      containsBegin = true;
    } else if (i === end) {
      containsEnd = true;
    }
    navigation.push({
      label: i.toString(),
      path: `${path}?${params.toString()}`,
    });
    count++;
  }
  if (!containsBegin) {
    params.set("page", "1");
    navigation = [{ label: "1", path: `${path}?${params.toString()}` }, { label: "...", path: "" }, ...navigation];
  }
  if (!containsEnd) {
    params.set("page", end.toString());
    navigation = [...navigation, { label: "...", path: "" }, { label: end.toString(), path: `${path}?${params.toString()}` }];
  }
</script>

<div class="flex flex-row items-center justify-center text-sm">
  {#each navigation as nav}
    {#if nav.label === "..."}
      <p class="m-1 flex h-10 w-10 select-none items-center justify-center rounded-md bg-white text-center font-medium shadow-sm shadow-gray-500">{nav.label}</p>
    {:else}
      <a class="m-1 flex h-10 w-10 items-center justify-center rounded-md bg-white text-center font-medium shadow-sm shadow-gray-500" data-sveltekit-preload-data="tap" href={nav.path}>{nav.label}</a>
    {/if}
  {/each}
</div>
