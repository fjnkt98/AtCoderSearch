<script lang="ts">
  import atcoderLogo from "$lib/assets/atcoder_logo.svg";
  import type { Problem } from "$lib/search";
  import dayjs from "dayjs";
  import timezone from "dayjs/plugin/timezone";
  import utc from "dayjs/plugin/utc";
  import { bgColorStyles } from "$lib/colors";

  dayjs.extend(timezone);
  dayjs.extend(utc);

  export let problem: Problem;

  function convertDateTime(date: string): string {
    return dayjs(date).tz("Asia/Tokyo").format("YYYY/MM/DD HH:mm:ss");
  }

  function convertDuration(duration: number): string {
    return `${duration / 60} minutes`;
  }
</script>

<div class="my-2 w-full min-w-[320px] max-w-5xl rounded-2xl bg-white px-2 py-2 shadow-md shadow-gray-300 sm:px-4 sm:py-4">
  <div class="flex flex-row items-center">
    <a href={problem.contest_url} target="_blank" rel="noreferrer">
      <img alt="AtCoder Logo" src={atcoderLogo} class="mx-2 my-1 aspect-square h-8 rounded-full bg-white sm:m-2 sm:h-12" />
    </a>
    <div class="mx-2">
      <a class="block break-all text-lg sm:my-1 sm:text-xl" href={problem.problem_url} target="_blank" rel="noreferrer">{problem.problem_title}</a>
      <a class="sm:text-md break-all text-sm text-blue-500" href={problem.problem_url} target="_blank" rel="noreferrer">{problem.problem_url}</a>
    </div>
  </div>

  <div class="my-2">
    <div class="m-1 px-2 text-sm text-slate-500">
      <a href={problem.contest_url} target="_blank" rel="noreferrer">
        {problem.contest_title}
      </a>
    </div>
    <div class="m-1 px-2 text-sm text-slate-500">
      {convertDateTime(problem.start_at)}
      {convertDuration(problem.duration)}
    </div>
  </div>

  <div class="mt-2 flex flex-row pt-1 text-xs font-bold text-white">
    <div class={`mx-1 rounded-full px-2 py-1 ${bgColorStyles.get(problem.category)}`}>
      {problem.category}
    </div>
    {#if problem.difficulty != null && problem.color != null}
      <div class={`mx-1 rounded-full px-2 py-1 ${bgColorStyles.get(problem.color)}`}>{problem.difficulty}</div>
    {/if}
  </div>
</div>
