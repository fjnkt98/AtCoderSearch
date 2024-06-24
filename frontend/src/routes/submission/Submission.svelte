<script lang="ts">
  import { textColorStyles } from "$lib/colors";
  import type { Submission } from "$lib/response";
  import dayjs from "dayjs";
  import timezone from "dayjs/plugin/timezone";
  import utc from "dayjs/plugin/utc";

  dayjs.extend(timezone);
  dayjs.extend(utc);

  function convertDateTime(date: string): string {
    return dayjs(date).tz("Asia/Tokyo").format("YYYY/MM/DD HH:mm:ss");
  }

  export let item: Submission;
  const colors = new Map<string, string>([
    ["AC", "bg-green-600 text-white"],
    ["WA", "bg-yellow-500 text-white"],
    ["TLE", "bg-red-600 text-white"],
    ["RE", "bg-red-600 text-white"],
    ["CE", "bg-red-600 text-white"],
    ["MLE", "bg-red-600 text-white"],
  ]);
</script>

<div class="mx-2 my-1 flex flex-row justify-between gap-1 rounded-lg bg-white px-3 py-2 shadow-sm shadow-gray-400 lg:min-w-96">
  <div class="flex flex-grow flex-col justify-around">
    <div class="flex flex-row items-center gap-2">
      <a class="block text-xs font-bold text-blue-600" href={item.submissionUrl} target="_blank" rel="noreferrer">#{item.submissionId}</a>
      <span class="block whitespace-nowrap text-xs text-gray-700">{convertDateTime(item.submittedAt)}</span>
    </div>
    <div class="my-0.5 flex flex-row items-center gap-2">
      <span class="block w-1/3 break-normal text-sm">{item.contestId.toUpperCase()}</span>
      <span class={`block flex-grow text-center text-sm font-medium ${textColorStyles.get(item.color ?? "black")}`}>{item.problemTitle}</span>
    </div>
  </div>

  <div class="flex flex-col items-end justify-between text-sm">
    <div class="flex flex-row items-center justify-between gap-1">
      <span class="block text-center text-sm">{item.userId}</span>
      <span class={`block rounded-xl px-1.5 py-0.5 text-xs ${colors.get(item.result) ?? "bg-gray-600 text-white"}`}>{item.result}</span>
    </div>
    <span class="block w-full text-nowrap text-center text-sm">{item.languageGroup}</span>

    <div class="flex flex-row items-center justify-between gap-2 text-xs">
      <span class="block text-nowrap text-center">{item.length} Byte</span>
      {#if item.executionTime != null}
        <span class="block text-nowrap text-center">{item.executionTime} ms</span>
      {/if}
    </div>
  </div>
</div>
