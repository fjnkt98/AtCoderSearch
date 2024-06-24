<script lang="ts">
  import { textColorStyles } from "$lib/colors";
  import type { User } from "$lib/response";
  import Icon from "svelte-icons-pack/Icon.svelte";
  import FaSolidCrown from "svelte-icons-pack/fa/FaSolidCrown";

  export let item: User;

  const crownColors = new Map<string, string>([
    ["crown_bronze", "#EF6C00"],
    ["crown_silver", "#B0BEC5"],
    ["crown_gold", "#FDD835"],
    ["crown_champion", "#F44336"],
  ]);
</script>

<div class="mx-2 my-1 flex flex-row rounded-lg bg-white px-1 py-2 shadow-sm shadow-gray-400 md:px-3 lg:min-w-96">
  <div class="flex flex-grow flex-row items-center justify-between">
    <div class="flex w-12 flex-col items-center justify-between">
      {#if item.crown != null}
        <Icon src={FaSolidCrown} color={crownColors.get(item.crown)} />
      {/if}
      <span class="block text-center">{item.rank}</span>
      {#if item.country != null}
        <span class="block text-center text-xs">{item.country}</span>
      {/if}
    </div>
    <div class="flex flex-grow flex-col">
      <div class="flex flex-row items-center justify-start gap-1">
        <a href={item.userUrl} class={`block flex-grow text-center ${textColorStyles.get(item.color ?? "black")}`} target="_blank">{item.userId}</a>
      </div>

      {#if item.affiliation != null}
        <span class="block text-pretty text-center text-xs text-gray-600">{item.affiliation}</span>
      {/if}
    </div>
  </div>

  <div class="flex-0 flex flex-col justify-center text-gray-900">
    <div class="flex flex-row items-center justify-between gap-4 text-xs">
      <span class="block text-nowrap">レーティング</span>
      <span class="block">{item.rating}</span>
    </div>
    {#if item.birthYear != null}
      <div class="flex flex-row items-center justify-between gap-4 text-xs">
        <span class="block">誕生年</span>
        <span class="block">{item.birthYear}</span>
      </div>
    {/if}
    <div class="flex flex-row items-center justify-between gap-4 text-xs">
      <span class="block">参加回数</span>
      <span class="block">{item.joinCount}</span>
    </div>
  </div>
</div>
