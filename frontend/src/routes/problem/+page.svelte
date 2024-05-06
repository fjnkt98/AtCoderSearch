<script lang="ts">
  import Tab from "$lib/Tab.svelte";
  import { page } from "$app/stores";
  import type { SearchProblemResult } from "$lib/response";
  import { intoURLSearchParams } from "$lib/request";
  import type { SearchProblemParameter } from "$lib/request";
  import Header from "$lib/Header.svelte";
  import dayjs from "dayjs";
  import timezone from "dayjs/plugin/timezone";
  import utc from "dayjs/plugin/utc";
  import SearchBar from "$lib/SearchBar.svelte";
  import { bgColorStyles, textColorStyles } from "$lib/colors";

  dayjs.extend(timezone);
  dayjs.extend(utc);

  let opened: boolean = false;
  let sort: string = "-startAt";
  let detail: string | null = null;

  const sorts = [
    ["-score", "検索スコア順"],
    ["startAt", "開催日時早い順"],
    ["-startAt", "開催日時遅い順"],
    ["difficulty", "難易度低い順"],
    ["-difficulty", "難易度高い順"],
  ];

  let param: SearchProblemParameter = {
    limit: 20,
    page: 1,
    q: $page.url.searchParams.get("q"),
    sort: ["-startAt"],
    facet: ["category", "difficulty"],
    category: null,
    difficultyFrom: null,
    difficultyTo: null,
    color: null,
    userId: null,
    difficulty: null,
    excludeSolved: null,
    experimental: null,
    prioritizeRecent: null,
  };

  const data: SearchProblemResult = {
    stats: {
      time: 99,
      total: 100,
      index: 1,
      pages: 5,
      count: 20,
      params: null,
      facet: null,
    },
    items: [
      {
        problemId: "agc035_f",
        problemTitle: "F. Two Histograms",
        problemUrl: "https://atcoder.jp/contests/agc035/tasks/agc035_f",
        contestId: "agc035",
        contestTitle: "AtCoder Grand Contest 035",
        contestUrl: "https://atcoder.jp/contests/agc035",
        difficulty: 3720,
        color: "gold",
        startAt: "2019-07-14T21:30:00+09:00",
        duration: 7800,
        rateChange: "All",
        category: "AGC",
        score: 1.0,
      },
      {
        problemId: "agc035_d",
        problemTitle: "D. Add and Remove",
        problemUrl: "https://atcoder.jp/contests/agc035/tasks/agc035_d",
        contestId: "agc035",
        contestTitle: "AtCoder Grand Contest 035",
        contestUrl: "https://atcoder.jp/contests/agc035",
        difficulty: 2902,
        color: "red",
        startAt: "2019-07-14T21:30:00+09:00",
        duration: 7800,
        rateChange: "All",
        category: "AGC",
        score: 1,
      },
    ],
    message: null,
  };

  function convertDateTime(date: string): string {
    return dayjs(date).tz("Asia/Tokyo").format("YYYY/MM/DD HH:mm:ss");
  }
</script>

<div class="flex h-dvh flex-col bg-gray-100">
  <Header />
  <Tab selected={"problem"} />

  <div class="my-2 flex flex-row items-center justify-between gap-2 px-2 text-sm">
    <div class={"w-1/5 flex-auto rounded-md border text-center align-middle shadow-sm " + (opened ? "border-gray-400 bg-green-600 font-semibold text-white" : "border-slate-500 bg-white")}>
      <label class="block h-full w-full cursor-pointer select-none py-1">
        絞り込み
        <input type="checkbox" bind:checked={opened} class="h-0 w-0" />
      </label>
    </div>
    <div class="flex-grow">
      <SearchBar href={"/problem"} />
    </div>
  </div>

  {#if opened}
    <div>ほげ</div>
  {:else}
    <div class="mx-2 mb-2 flex flex-row justify-between text-sm">
      <select
        class="block rounded-md border border-gray-400 bg-white px-2 py-1"
        bind:value={sort}
        on:change={() => {
          param = { ...param, sort: [sort] };
        }}
      >
        {#each sorts as s}
          <option value={s[0]}>
            {s[1]}
          </option>
        {/each}
      </select>
      <p class="w-full text-right text-gray-700">{data.stats.count}件/{data.stats.total}件 約{data.stats.time}ms</p>
    </div>

    {#each data.items as item}
      <div class="mx-2 my-0.5 flex flex-col rounded-lg bg-white px-3 py-3 shadow-md shadow-gray-200">
        <div class="flex flex-row items-center justify-between">
          <a href={item.contestUrl} target="_blank" class="block font-medium text-blue-500">
            {item.contestId.toUpperCase()}
          </a>
          <a href={item.problemUrl} target="_blank" class={`block text-pretty font-medium ${textColorStyles.get(item.color ?? "black")}`}>
            {item.problemTitle}
          </a>
          <label class={`text-sm text-black ` + (detail == item.problemId ? "" : "underline")}>
            {detail == item.problemId ? "閉じる" : "詳細"}
            <input
              type="button"
              class="h-0 w-0"
              on:click={() => {
                detail = detail == item.problemId ? null : item.problemId;
              }}
            />
          </label>
        </div>
        {#if detail == item.problemId}
          <div class="py-1">
            <div class="my-0.5 flex flex-row items-end justify-between">
              <span class="block text-sm text-slate-500">{item.contestTitle}</span>
              <span class={`block rounded-full px-2 py-0.5 text-xs font-medium text-white  ${bgColorStyles.get(item.category)}`}>{item.category}</span>
            </div>
            <div class="flex flex-row items-center justify-between">
              <span class="block text-sm text-slate-500">{convertDateTime(item.startAt)}</span>
              {#if item.difficulty != null}
                <span class={`block rounded-full px-2 py-0.5 text-xs font-medium text-white ${bgColorStyles.get(item.color ?? "black")}`}>{item.difficulty}</span>
              {/if}
            </div>
          </div>
        {/if}
      </div>
    {/each}
  {/if}
</div>
