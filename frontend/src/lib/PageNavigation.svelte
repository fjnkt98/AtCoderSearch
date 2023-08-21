<script lang="ts">
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  export let end: number;
  export let current: number;

  const width = 5;

  const generateLabels = (c: number, e: number): string[] => {
    let labels: string[] = [];
    let containsBegin = false;
    let containsEnd = false;
    for (let i = Math.max(c - width, 1); i <= Math.min(c + width, e); i++) {
      if (i === 1) {
        containsBegin = true;
      } else if (i === e) {
        containsEnd = true;
      }
      labels.push(i.toString());
    }
    if (!containsBegin) {
      labels = ["1", "...", ...labels];
    }
    if (!containsEnd) {
      labels = [...labels, "...", e.toString()];
    }
    return labels;
  };

  $: labels = generateLabels(current, end);
</script>

<div class="flex flex-row items-center justify-center text-sm">
  {#each labels as label}
    {#if label === "..."}
      <p class="m-1 flex h-10 w-10 select-none items-center justify-center rounded-md bg-white text-center font-medium shadow-sm shadow-gray-500">{label}</p>
    {:else}
      <button
        class="m-1 flex h-10 w-10 items-center justify-center rounded-md bg-white text-center font-medium shadow-sm shadow-gray-500"
        on:click={() => {
          const params = new URLSearchParams($page.url.searchParams);
          params.set("page", label);
          goto(`${$page.url.pathname}?${params.toString()}`);
        }}
        >{label}
      </button>
    {/if}
  {/each}
</div>
