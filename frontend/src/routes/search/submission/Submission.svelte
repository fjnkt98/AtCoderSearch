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
</script>

<div class="my-2 w-2/3 min-w-[600px] rounded-2xl bg-white px-4 py-2 shadow-md shadow-gray-300">
  <div class="flex flex-row items-center">
    <div class="mx-2 flex flex-grow flex-col pr-4">
      <a class="inline-block text-sm" href={submission.submission_url} target="_blank" rel="noreferrer">#{submission.submission_id}</a>
      <div class="flex flex-row items-center justify-between">
        <span class="inline-block text-lg">{submission.problem_title}</span>
        <span class="mt-2 inline-block text-lg">{submission.user_id}</span>
      </div>
      <div class="flex flex-row items-center justify-between">
        <span class="inline-block break-keep text-sm text-gray-500">{submission.contest_title}</span>
        <span class="inline-block whitespace-nowrap text-sm text-gray-500">{convertDateTime(submission.submitted_at)}</span>
      </div>
    </div>

    <div class="flex flex-col items-center justify-between">
      <span class="inline-block whitespace-nowrap text-sm">{submission.language}</span>
      <span class="inline-block text-sm">{submission.length} Byte</span>
      {#if submission.execution_time != null}
        <span class="inline-block text-sm">{submission.execution_time} ms</span>
      {/if}
      <span class="inline-block text-lg">{submission.result}</span>
    </div>
  </div>
</div>
