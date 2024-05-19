<script lang="ts">
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
  /**
   *   submissionId: number;
  submittedAt: string;
  submissionUrl: string;
  problemId: string;
  problemTitle: string;
  contestId: string;
  contestTitle: string;
  category: string;
  difficulty: number;
  color: string;
  userId: string;
  language: string;
  point: number;
  length: number;
  result: string;
  executionTime: number | null;
   * 
  */
  const colors = new Map<string, string>([
    ["AC", "bg-green-600 text-white"],
    ["WA", "bg-yellow-500 text-white"],
    ["TLE", "bg-red-600 text-white"],
    ["RE", "bg-red-600 text-white"],
    ["CE", "bg-red-600 text-white"],
    ["MLE", "bg-red-600 text-white"],
  ]);
</script>

<div class="mx-2 my-1 flex min-w-96 flex-row justify-between rounded-lg bg-white px-3 py-2 shadow-sm shadow-gray-400">
  <div class="flex flex-col justify-between">
    <div class="flex flex-row items-center gap-3">
      <a class="block font-bold text-blue-600" href={item.submissionUrl} target="_blank" rel="noreferrer">#{item.submissionId}</a>
      <span class="block whitespace-nowrap text-sm text-gray-700">{convertDateTime(item.submittedAt)}</span>
    </div>
    <div class="flex flex-row items-center gap-3">
      <span class="block break-all text-sm">{item.contestId.toUpperCase()}</span>
      <span class="block break-all">{item.problemTitle}</span>
    </div>
  </div>

  <div class="flex min-w-32 flex-col items-end justify-between">
    <div class="flex flex-row items-center justify-between gap-3">
      <span class="block break-all">{item.userId}</span>
      <span class={`px- block rounded-xl px-2 py-0.5 text-xs ${colors.get(item.result) ?? "bg-gray-600 text-white"}`}>{item.result}</span>
    </div>
    <span class="block whitespace-nowrap text-xs">{item.language}</span>

    <div class="flex flex-row items-center justify-between gap-3 text-xs">
      <span class="block">{item.length} Byte</span>
      {#if item.executionTime != null}
        <span class="block">{item.executionTime} ms</span>
      {/if}
    </div>
  </div>
</div>
