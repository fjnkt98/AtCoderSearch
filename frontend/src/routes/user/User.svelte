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

<div class="mx-2 my-1 flex min-w-96 flex-row rounded-lg bg-white px-3 py-2 shadow-sm shadow-gray-400">
  <div class="flex flex-grow flex-row items-center justify-between">
    <span class="block w-16 text-center">{item.rank}</span>
    <div class="flex flex-grow flex-col">
      <div class="flex flex-row items-center justify-start gap-1">
        {#if item.crown != null}
          <Icon src={FaSolidCrown} color={crownColors.get(item.crown)} />
        {/if}
        <a href={item.userUrl} class={`block flex-grow  text-xl ${textColorStyles.get(item.color ?? "black")}`} target="_blank">{item.userId}</a>
      </div>

      {#if item.affiliation != null}
        <span class="block text-sm text-gray-600">{item.affiliation}</span>
      {/if}
    </div>
  </div>

  <div class="flex-0 flex min-w-32 flex-col justify-center text-gray-900">
    <div class="flex flex-row items-center justify-between gap-4 text-xs">
      <span class="block">レーティング</span>
      <span class="block">{item.rating}</span>
    </div>

    {#if item.country != null}
      <div class="flex flex-row items-center justify-between gap-4 text-xs">
        <span class="block">国と地域</span>
        <span class="block">{item.country}</span>
      </div>
    {/if}
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
