<script lang="ts">
  import { defaultProblemSearchParams, defaultUserSearchParams, defaultSubmissionSearchParams, defaultRecommendParams } from "$lib/defaults";
</script>

<nav class="flex w-full items-center justify-start border">
  <a href="/" class="font-roboto mx-4 my-2 text-3xl">AtCoder Search</a>
</nav>

<div class="flex flex-1 justify-center overflow-auto p-2 py-8">
  <div class="grid h-4/6 w-5/6 max-w-6xl gap-4 p-2 sm:grid-cols-1 md:grid-cols-2">
    <a class="block rounded-3xl bg-white p-4 no-underline shadow-lg" href={`/search/problem?${defaultProblemSearchParams.toString()}`}>
      <h3 class="font-roboto text-2xl font-medium">Search Problems</h3>
      <p class="py-3 text-gray-600">AtCoderの問題を検索します</p>
    </a>
    <a class="block rounded-3xl bg-white p-4 no-underline shadow-lg" href={`/search/user?${defaultUserSearchParams.toString()}`}>
      <h3 class="font-roboto text-2xl font-medium">Search Users</h3>
      <p class="py-3 text-gray-600">AtCoderのユーザを検索します</p>
    </a>
    <a class="block rounded-3xl bg-white p-4 no-underline shadow-lg" href={`/search/submission?${defaultSubmissionSearchParams.toString()}`}>
      <h3 class="font-roboto text-2xl font-medium">Submissions</h3>
      <p class="py-3 text-gray-600">1ヶ月以内の提出を表示します</p>
    </a>
    <a class="block rounded-3xl bg-white p-4 no-underline shadow-lg" href={`/recommend/problem?${defaultRecommendParams.toString()}`}>
      <h3 class="font-roboto text-2xl font-medium">Recommend <span class="text-red-500">(alpha)</span></h3>
      <p class="py-3 text-gray-600">おすすめの問題を提示します</p>
    </a>
  </div>
</div>
