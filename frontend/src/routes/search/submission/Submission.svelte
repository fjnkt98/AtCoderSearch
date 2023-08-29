<script lang="ts">
  import type { Submission } from "$lib/search";
  import dayjs from "dayjs";
  import timezone from "dayjs/plugin/timezone";
  import utc from "dayjs/plugin/utc";

  dayjs.extend(timezone);
  dayjs.extend(utc);

  export let submission: Submission;

  function convertDateTime(date: string): string {
    return dayjs(date).tz("Asia/Tokyo").format("YYYY/MM/DD HH:mm:ss");
  }

  const colors = new Map<string, string>([
    ["AC", "bg-green-600 text-white"],
    ["WA", "bg-yellow-500 text-white"],
    ["TLE", "bg-red-600 text-white"],
    ["RE", "bg-red-600 text-white"],
    ["CE", "bg-red-600 text-white"],
    ["MLE", "bg-red-600 text-white"],
  ]);
</script>

<div class="my-2 w-2/3 min-w-[600px] rounded-2xl bg-white px-4 py-2 shadow-md shadow-gray-300">
  <div class="flex flex-row items-center justify-between">
    <div class="mx-2 flex basis-4/5 flex-col pr-4">
      <a class="inline-block text-sm font-bold text-blue-600" href={submission.submission_url} target="_blank" rel="noreferrer">#{submission.submission_id}</a>
      <div class="flex flex-row items-center justify-between">
        <span class="inline-block text-lg">{submission.problem_title}</span>
        <span class="mt-2 inline-block text-lg">{submission.user_id}</span>
      </div>
      <div class="flex flex-row items-end justify-between">
        <span class="inline-block break-all text-sm text-gray-500">{submission.contest_title}</span>
        <span class="inline-block whitespace-nowrap text-sm text-gray-500">{convertDateTime(submission.submitted_at)}</span>
      </div>
    </div>

    <div class="flex basis-1/5 flex-col items-center justify-between">
      <span class="inline-block whitespace-nowrap text-sm">{submission.language}</span>
      <span class="inline-block text-sm">{submission.length} Byte</span>
      {#if submission.execution_time != null}
        <span class="inline-block text-sm">{submission.execution_time} ms</span>
      {/if}
      <span class={`inline-block rounded-xl px-3 text-lg ${colors.get(submission.result) ?? "bg-gray-600 text-white"}`}>{submission.result}</span>
    </div>
  </div>
</div>
