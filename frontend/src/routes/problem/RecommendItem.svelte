<script lang="ts">
  import { bgColorStyles, textColorStyles } from "$lib/colors";
  import type { Problem } from "$lib/response";

  export let item: Problem;
</script>

<div class="flex w-48 flex-none flex-col text-nowrap rounded-md px-3 py-2 shadow-md shadow-gray-400">
  <div class="flex flex-row items-center justify-between">
    <a href={item.contestUrl} class="block text-nowrap text-xs text-gray-700" target="_blank">{item.contestId.toUpperCase()}</a>
    {#if item.difficulty != null}
      <div class="flex flex-row">
        <span class={`rounded-full ${bgColorStyles.get(item.color ?? "black")} px-2 py-0.5 text-xs text-white`}>{item.difficulty}</span>
      </div>
    {/if}
  </div>
  <a href={item.problemUrl} class={`${textColorStyles.get(item.color ?? "black")} mt-2 block flex-grow text-wrap text-sm`} target="_blank">{item.problemTitle}</a>
</div>
