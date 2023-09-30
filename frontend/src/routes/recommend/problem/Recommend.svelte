<script lang="ts">
  import atcoderLogo from "$lib/assets/atcoder_logo.svg";
  import type { Recommend, Problem } from "$lib/search";
  import dayjs from "dayjs";
  import timezone from "dayjs/plugin/timezone";
  import utc from "dayjs/plugin/utc";
  import { bgColorStyles } from "$lib/colors";

  dayjs.extend(timezone);
  dayjs.extend(utc);

  export let problem: Recommend | Problem;

  function convertDateTime(date: string): string {
    return dayjs(date).tz("Asia/Tokyo").format("YYYY/MM/DD HH:mm:ss");
  }
</script>

<div class="mx-2 my-4 flex min-w-[320px] flex-col justify-between rounded-2xl bg-white px-4 py-4 shadow-md shadow-gray-300">
  <div class="flex flex-row items-center">
    <a href={problem.contest_url} target="_blank" rel="noreferrer">
      <img alt="AtCoder Logo" src={atcoderLogo} class="m-1 aspect-square h-12 rounded-full bg-white" />
    </a>
    <div class="mx-2">
      <a class="block text-xl" href={problem.problem_url} target="_blank" rel="noreferrer">{problem.problem_title}</a>
    </div>
  </div>

  <div>
    <div class="px-2 text-sm text-gray-600">
      <a href={problem.contest_url} target="_blank" rel="noreferrer">
        {problem.contest_title}
      </a>
    </div>
    <div class="m-1 px-2 text-sm text-slate-500">
      {convertDateTime(problem.start_at)}
    </div>
  </div>

  <div class="flex flex-row justify-between pt-1 text-xs">
    <div class="flex flex-row">
      <div class={`mx-1 rounded-full px-2 py-1 font-bold text-white ${bgColorStyles.get(problem.category)}`}>
        {problem.category}
      </div>
      {#if problem.difficulty != null && problem.color != null}
        <div class={`mx-1 rounded-full px-2 py-1 font-bold text-white ${bgColorStyles.get(problem.color)}`}>{problem.difficulty}</div>
      {/if}
    </div>
    {#if "score" in problem}
      <div class="mx-1 rounded-full px-2 py-1 text-gray-500">{problem.score}</div>
    {/if}
  </div>
</div>
