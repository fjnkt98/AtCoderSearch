<script lang="ts">
  import { _ } from "$env/static/private";
  import atcoderLogo from "$lib/assets/atcoder_logo.svg";
  import type { ProblemSearchResult } from "$lib/search";
  import dayjs from "dayjs";
  import timezone from "dayjs/plugin/timezone";
  import utc from "dayjs/plugin/utc";

  dayjs.extend(timezone);
  dayjs.extend(utc);

  export let data: ProblemSearchResult;

  function convertDateTime(date: string): string {
    return dayjs(date).tz("Asia/Tokyo").format("YYYY/MM/DD HH:mm:ss");
  }

  function convertDuration(duration: number): string {
    return `${duration / 60} minutes`;
  }

  const colors = new Map<string, string>([
    ["ABC", "bg-blue-500"],
    ["ABC-Like", "bg-sky-600"],
    ["AGC", "bg-yellow-600"],
    ["AGC-Like", "bg-amber-500"],
    ["AHC", "bg-green-500"],
    ["ARC", "bg-red-500"],
    ["ARC-Like", "bg-orange-700"],
    ["JAG", "bg-slate-500"],
    ["JOI", "bg-slate-600"],
    ["Marathon", "bg-slate-600"],
    ["Other Contests", "bg-slate-600"],
    ["Other Sponsored", "bg-slate-600"],
    ["PAST", "bg-slate-600"],
    ["black", "bg-gray-900"],
    ["gray", "bg-gray-500"],
    ["brown", "bg-amber-800"],
    ["green", "bg-green-500"],
    ["cyan", "bg-cyan-400"],
    ["blue", "bg-blue-600"],
    ["yellow", "bg-yellow-300"],
    ["orange", "bg-orange-400"],
    ["red", "bg-red-600"],
    ["silver", "bg-zinc-400"],
    ["gold", "bg-amber-400"],
  ]);
</script>

<div class="flex-1 overflow-auto p-16">
  {#await data.items}
    <p>search...</p>
  {:then items}
    {#each items as item}
      <div class="my-2 rounded-2xl bg-white px-4 py-4 shadow-md shadow-gray-300">
        <div class="flex flex-row items-center">
          <a href={item.contest_url} target="_blank" rel="noreferrer">
            <img alt="AtCoder Logo" src={atcoderLogo} class="m-2 aspect-square h-12 rounded-full bg-white" />
          </a>
          <div class="mx-2">
            <p class="my-1 text-xl">{item.problem_title}</p>
            <a class="text-md text-blue-500" href={item.problem_url} target="_blank" rel="noreferrer">{item.problem_url}</a>
          </div>
        </div>

        <div class="my-2">
          <div class="m-1 px-2 text-sm text-slate-500">{item.contest_title}</div>
          <div class="m-1 px-2 text-sm text-slate-500">
            {convertDateTime(item.start_at)}
            {convertDuration(item.duration)}
          </div>
        </div>

        <div class="mt-2 flex flex-row pt-1 text-xs font-bold text-white">
          <div class={`mx-1 rounded-full px-2 py-1 ${colors.get(item.category)}`}>
            {item.category}
          </div>
          {#if item.difficulty != null && item.color != null}
            <div class={`mx-1 rounded-full px-2 py-1 ${colors.get(item.color)}`}>{item.difficulty}</div>
          {/if}
        </div>
      </div>
    {/each}
  {:catch}
    <p>Search failed!</p>
  {/await}
</div>
