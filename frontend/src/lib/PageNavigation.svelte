<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";

  export let end: number = 999;
  export let current: number = 1;
  export let enableEnd: boolean = true;

  const width = 4;

  const generateLabels = (c: number, e: number): string[] => {
    if ((c === 1 && e === 1) || e === 0) {
      return ["1"];
    }
    let labels: string[] = [];
    let containsBegin = false;
    let containsEnd = false;
    const left = Math.max(c - Math.floor(width / 2), 1);
    const right = Math.min(left + width, e);
    for (let i = left; i <= right; i++) {
      if (i === 1) {
        containsBegin = true;
      } else if (i === e) {
        containsEnd = true;
      }
      labels.push(i.toString());
    }
    if (!containsBegin) {
      if (labels[0] == "2") {
        labels = ["1", ...labels];
      } else {
        labels = ["1", "...", ...labels];
      }
    }
    if (!containsEnd) {
      if (enableEnd) {
        labels = [...labels, "...", e.toString()];
      } else {
        labels = [...labels, "..."];
      }
    }
    return labels;
  };

  $: labels = generateLabels(current, end);
</script>

<div class="flex flex-row items-center justify-center text-sm">
  {#each labels as label}
    {#if label === "..."}
      <p class="m-1 flex h-8 w-8 select-none items-center justify-center rounded-md bg-white text-center font-medium shadow-sm shadow-gray-500 sm:h-10 sm:w-10">{label}</p>
    {:else}
      <button
        class={"m-1 flex h-8 w-8 items-center justify-center rounded-md text-center font-medium shadow-sm shadow-gray-500 hover:bg-gray-100 sm:h-10 sm:w-10 " +
          (current.toString() === label ? "bg-gray-700 text-white hover:bg-gray-800" : "bg-white")}
        on:click={() => {
          const params = new URLSearchParams($page.url.searchParams);
          params.set("p", label);
          goto(`${$page.url.pathname}?${params.toString()}`);
        }}
        >{label}
      </button>
    {/if}
  {/each}
</div>
