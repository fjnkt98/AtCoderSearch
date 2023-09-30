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

<div class="my-2 w-full min-w-[320px] max-w-5xl rounded-2xl bg-white px-2 py-2 shadow-md shadow-gray-300 sm:px-4 sm:py-3">
  <div class="flex flex-row items-center justify-between">
    <div class="mx-2 flex basis-4/5 flex-col pr-4">
      <a class="inline-block text-sm font-bold text-blue-600" href={submission.submission_url} target="_blank" rel="noreferrer">#{submission.submission_id}</a>
      <span class="inline-block break-all py-1 text-base sm:text-lg">{submission.problem_title}</span>
      <div class="flex flex-col items-start justify-between">
        <span class="inline-block break-all text-xs text-gray-500 sm:text-sm">{submission.contest_title}</span>
        <span class="inline-block whitespace-nowrap text-sm text-gray-500 sm:text-sm">{convertDateTime(submission.submitted_at)}</span>
      </div>
    </div>

    <div class="flex basis-1/5 flex-col items-center justify-between">
      <span class="inline-block break-all text-base sm:text-lg">{submission.user_id}</span>
      <span class="inline-block whitespace-nowrap text-xs sm:text-sm">{submission.language}</span>
      <span class="inline-block text-xs sm:text-sm">{submission.length} Byte</span>
      {#if submission.execution_time != null}
        <span class="inline-block text-xs sm:text-sm">{submission.execution_time} ms</span>
      {/if}
      <span class={`inline-block rounded-xl px-2 py-0.5 text-base sm:px-4 sm:text-base ${colors.get(submission.result) ?? "bg-gray-600 text-white"}`}>{submission.result}</span>
    </div>
  </div>
</div>
